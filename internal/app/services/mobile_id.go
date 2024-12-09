package services

import (
	"context"
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
	cfg   *config.Config
	log   *logger.Logger
	debug bool
}

func NewMobileId(cfg *config.Config, log *logger.Logger) MobileIdProvider {
	return &mobileIdProvider{
		cfg:   cfg,
		debug: cfg.LogLevel == config.DebugLevel,
		log:   log,
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
	if s.debug {
		client.EnableTrace()
	}

	response, err := client.R().
		SetHeader("Accept", "application/json").
		SetHeader("Content-Type", "application/json").
		SetBody(body).
		Post(endpoint)
	if s.debug {
		debug(response, err)
	}
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
	if s.debug {
		client.EnableTrace()
	}

	response, err := client.R().
		SetQueryParam("timeoutMs", MobileIdTimeout).
		SetHeader("Accept", "application/json").
		SetHeader("Content-Type", "application/json").
		Get(endpoint)
	if s.debug {
		debug(response, err)
	}
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
