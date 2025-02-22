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

const (
	MobileIdTextFormat = "GSM-7"
	MobileIdHashType   = "SHA512"
	MobileIdTimeout    = "120000"
)

type MobileIdRequestBody struct {
	RelyingPartyName       string `json:"relyingPartyName"`
	RelyingPartyUUID       string `json:"relyingPartyUUID"`
	NationalIdentityNumber string `json:"nationalIdentityNumber"`
	PhoneNumber            string `json:"phoneNumber"`
	Hash                   string `json:"hash"`
	HashType               string `json:"hashType"`
	Language               string `json:"language"`
	DisplayText            string `json:"displayText"`
	DisplayTextFormat      string `json:"displayTextFormat"`
}

type MobileIdProvider interface {
	CreateSession(ctx context.Context, params dto.CreateMobileIdSessionRequest) (*dto.MobileIdProviderSessionResponse, error)
	GetSessionStatus(id uuid.UUID) (*dto.MobileIdProviderSessionStatusResponse, error)
}

type mobileIdProvider struct {
	cfg *config.Config
	log *logger.Logger
}

func NewMobileId(cfg *config.Config, log *logger.Logger) MobileIdProvider {
	return &mobileIdProvider{
		cfg: cfg,
		log: log,
	}
}

func (s *mobileIdProvider) CreateSession(_ context.Context, params dto.CreateMobileIdSessionRequest) (*dto.MobileIdProviderSessionResponse, error) {
	hash, err := generateHash()
	if err != nil {
		return nil, err
	}

	endpoint := fmt.Sprintf("%s/authentication", s.cfg.MobileId.BaseURL)

	body := MobileIdRequestBody{
		RelyingPartyName:       s.cfg.MobileId.RelyingPartyName,
		RelyingPartyUUID:       s.cfg.MobileId.RelyingPartyUUID,
		PhoneNumber:            params.PhoneNumber,
		NationalIdentityNumber: params.PersonalCode,
		Hash:                   hash,
		HashType:               MobileIdHashType,
		Language:               params.Locale,
		DisplayText:            s.cfg.MobileId.Text,
		DisplayTextFormat:      MobileIdTextFormat,
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
		return nil, errors.ErrMobileIdProviderError
	}

	var result dto.MobileIdProviderSessionResponse
	if err = json.Unmarshal(response.Body(), &result); err != nil {
		return nil, err
	}

	code, err := generateCode(hash)
	if err != nil {
		return nil, err
	}

	return &dto.MobileIdProviderSessionResponse{
		ID:   result.ID,
		Code: code,
	}, nil
}

func (s *mobileIdProvider) GetSessionStatus(id uuid.UUID) (*dto.MobileIdProviderSessionStatusResponse, error) {
	endpoint := fmt.Sprintf("%s/authentication/session/%s", s.cfg.MobileId.BaseURL, id)

	client := resty.New()
	response, err := client.R().
		SetQueryParam("timeoutMs", MobileIdTimeout).
		SetHeader("Accept", "application/json").
		SetHeader("Content-Type", "application/json").
		Get(endpoint)
	if err != nil {
		return nil, err
	}

	if response.StatusCode() != 200 {
		return nil, errors.ErrMobileIdProviderError
	}

	var result dto.MobileIdProviderSessionStatusResponse
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
