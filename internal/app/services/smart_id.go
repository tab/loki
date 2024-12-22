package services

import (
	"context"
	"crypto/rand"
	"crypto/sha256"
	"crypto/sha512"
	"encoding/base64"
	"encoding/binary"
	"encoding/json"
	"fmt"

	"github.com/go-resty/resty/v2"
	"github.com/google/uuid"

	"loki/internal/app/errors"
	"loki/internal/app/models/dto"
	"loki/internal/config"
	"loki/pkg/logger"
)

type SmartIdAllowedInteraction struct {
	Type          string `json:"type"`
	DisplayText60 string `json:"displayText60"`
}

type SmartIdRequestBody struct {
	RelyingPartyName                string                      `json:"relyingPartyName"`
	RelyingPartyUUID                string                      `json:"relyingPartyUUID"`
	NationalIdentityNumber          string                      `json:"nationalIdentityNumber"`
	CertificateLevel                string                      `json:"certificateLevel"`
	SmartIdAllowedInteractionsOrder []SmartIdAllowedInteraction `json:"allowedInteractionsOrder"`
	Hash                            string                      `json:"hash"`
	HashType                        string                      `json:"hashType"`
}

const (
	SmartIdCertificateLevel = "QUALIFIED"
	SmartIdHashType         = "SHA512"
	SmartIdInteractionType  = "displayTextAndPIN"
	SmartIdTimeout          = "120000"
)

type SmartIdProvider interface {
	CreateSession(ctx context.Context, params dto.CreateSmartIdSessionRequest) (*dto.SmartIdProviderSessionResponse, error)
	GetSessionStatus(id uuid.UUID) (*dto.SmartIdProviderSessionStatusResponse, error)
}

type smartIdProvider struct {
	cfg *config.Config
	log *logger.Logger
}

func NewSmartId(cfg *config.Config, log *logger.Logger) SmartIdProvider {
	return &smartIdProvider{
		cfg: cfg,
		log: log,
	}
}

func (s *smartIdProvider) CreateSession(_ context.Context, params dto.CreateSmartIdSessionRequest) (*dto.SmartIdProviderSessionResponse, error) {
	hash, err := generateHash()
	if err != nil {
		return nil, err
	}

	nationalIdentityNumber := fmt.Sprintf("PNO%s-%s", params.Country, params.PersonalCode)
	endpoint := fmt.Sprintf("%s/authentication/etsi/%s", s.cfg.SmartId.BaseURL, nationalIdentityNumber)

	body := SmartIdRequestBody{
		RelyingPartyName:       s.cfg.SmartId.RelyingPartyName,
		RelyingPartyUUID:       s.cfg.SmartId.RelyingPartyUUID,
		NationalIdentityNumber: nationalIdentityNumber,
		CertificateLevel:       SmartIdCertificateLevel,
		Hash:                   hash,
		HashType:               SmartIdHashType,
		SmartIdAllowedInteractionsOrder: []SmartIdAllowedInteraction{
			{
				Type:          SmartIdInteractionType,
				DisplayText60: s.cfg.SmartId.Text,
			},
		},
	}

	client := resty.New()
	response, err := client.R().
		SetHeader("Accept", "application/json").
		SetHeader("Content-Type", "application/json").
		SetBody(body).
		Post(endpoint)
	if err != nil {
		return nil, err
	}

	if response.StatusCode() != 200 {
		return nil, errors.ErrSmartIdProviderError
	}

	var result dto.SmartIdProviderSessionResponse
	if err = json.Unmarshal(response.Body(), &result); err != nil {
		return nil, err
	}

	code, err := generateCode(hash)
	if err != nil {
		return nil, err
	}

	return &dto.SmartIdProviderSessionResponse{
		ID:   result.ID,
		Code: code,
	}, nil
}

func (s *smartIdProvider) GetSessionStatus(id uuid.UUID) (*dto.SmartIdProviderSessionStatusResponse, error) {
	endpoint := fmt.Sprintf("%s/session/%s", s.cfg.SmartId.BaseURL, id)

	client := resty.New()
	response, err := client.R().
		SetQueryParam("timeoutMs", SmartIdTimeout).
		SetHeader("Accept", "application/json").
		SetHeader("Content-Type", "application/json").
		Get(endpoint)
	if err != nil {
		return nil, err
	}

	if response.StatusCode() != 200 {
		return nil, errors.ErrSmartIdProviderError
	}

	var result dto.SmartIdProviderSessionStatusResponse
	if err = json.Unmarshal(response.Body(), &result); err != nil {
		return nil, err
	}

	return &result, nil
}

func generateHash() (string, error) {
	randBytes := make([]byte, 64)
	_, err := rand.Read(randBytes)
	if err != nil {
		return "", err
	}

	hash := sha512.Sum512(randBytes)
	encoded := base64.StdEncoding.EncodeToString(hash[:])

	return encoded, nil
}

func generateCode(hash string) (string, error) {
	decodedHash, err := base64.StdEncoding.DecodeString(hash)
	if err != nil {
		return "", err
	}

	sha256Hash := sha256.Sum256(decodedHash)
	lastTwoBytes := sha256Hash[len(sha256Hash)-2:]
	codeInt := binary.BigEndian.Uint16(lastTwoBytes)
	vc := codeInt % 10000
	code := fmt.Sprintf("%04d", vc)

	return code, nil
}
