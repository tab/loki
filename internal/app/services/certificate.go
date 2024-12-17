package services

import (
	"crypto/x509"
	"encoding/base64"
	"regexp"
	"strings"

	"loki/internal/app/errors"
	"loki/pkg/logger"
)

type CertificatePayload struct {
	IdentityNumber string
	PersonalCode   string
	FirstName      string
	LastName       string
}

type Certificate interface {
	Extract(value string) (*CertificatePayload, error)
}

type certificate struct {
	log *logger.Logger
}

func NewCertificate(log *logger.Logger) Certificate {
	return &certificate{
		log: log,
	}
}

func (c *certificate) Extract(value string) (*CertificatePayload, error) {
	certBytes, err := base64.StdEncoding.DecodeString(value)
	if err != nil {
		c.log.Error().Err(err).Msg("failed to decode certificate")
		return nil, err
	}

	cert, err := x509.ParseCertificate(certBytes)
	if err != nil {
		c.log.Error().Err(err).Msg("failed to parse certificate")
		return nil, err
	}

	subject := cert.Subject
	commonName := subject.CommonName

	parts := strings.Split(commonName, ",")
	if len(parts) < 2 {
		c.log.Error().Msgf("invalid CommonName format: %s", commonName)
		return nil, errors.ErrInvalidCertificate
	}

	personalCode, _ := extractPersonalCode(subject.SerialNumber)
	firstName := strings.TrimSpace(parts[0])
	lastName := strings.TrimSpace(parts[1])

	return &CertificatePayload{
		IdentityNumber: subject.SerialNumber,
		PersonalCode:   personalCode,
		FirstName:      firstName,
		LastName:       lastName,
	}, nil
}

func extractPersonalCode(identityNumber string) (string, error) {
	const prefix = "PNO"

	if !strings.HasPrefix(identityNumber, prefix) {
		return "", errors.ErrInvalidIdentityNumber
	}

	re := regexp.MustCompile(`PNO[A-Z]{2}-(\d+)`)
	matches := re.FindStringSubmatch(identityNumber)

	if len(matches) != 2 {
		return "", errors.ErrInvalidIdentityNumber
	}

	return matches[1], nil
}
