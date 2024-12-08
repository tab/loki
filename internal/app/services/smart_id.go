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

	"loki/internal/app/models/dto"
	"loki/internal/app/repositories"
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
	CertificateLevel = "QUALIFIED"
	HashType         = "SHA512"
	InteractionType  = "displayTextAndPIN"
	Timeout          = "120000"
)

type SmartIdProvider interface {
	CreateSession(ctx context.Context, params dto.CreateSmartIdSessionRequest) (*dto.SmartIdProviderSessionResponse, error)
	GetSessionStatus(id uuid.UUID) (*dto.SmartIdProviderSessionStatusResponse, error)
}

type smartIdProvider struct {
	cfg   *config.Config
	repo  repositories.Database
	log   *logger.Logger
	debug bool
}

func NewSmartId(
	cfg *config.Config,
	repo repositories.Database,
	log *logger.Logger,
) SmartIdProvider {

	return &smartIdProvider{
		cfg:   cfg,
		repo:  repo,
		log:   log,
		debug: cfg.LogLevel == config.DebugLevel,
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
		CertificateLevel:       CertificateLevel,
		Hash:                   hash,
		HashType:               HashType,
		SmartIdAllowedInteractionsOrder: []SmartIdAllowedInteraction{
			{
				Type:          InteractionType,
				DisplayText60: s.cfg.SmartId.Text,
			},
		},
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
	if s.debug {
		client.EnableTrace()
	}

	response, err := client.R().
		SetQueryParam("timeoutMs", Timeout).
		SetHeader("Accept", "application/json").
		SetHeader("Content-Type", "application/json").
		Get(endpoint)

	if s.debug {
		debug(response, err)
	}

	if err != nil {
		return nil, err
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

func debug(response *resty.Response, err error) {
	fmt.Println("Response Info:")
	fmt.Println("  Error      :", err)
	fmt.Println("  Status Code:", response.StatusCode())
	fmt.Println("  Status     :", response.Status())
	fmt.Println("  Proto      :", response.Proto())
	fmt.Println("  Time       :", response.Time())
	fmt.Println("  Received At:", response.ReceivedAt())
	fmt.Println("  Body       :\n", response.Body())
	fmt.Println()

	fmt.Println("Request Trace Info:")
	ti := response.Request.TraceInfo()
	fmt.Println("  DNSLookup     :", ti.DNSLookup)
	fmt.Println("  ConnTime      :", ti.ConnTime)
	fmt.Println("  TCPConnTime   :", ti.TCPConnTime)
	fmt.Println("  TLSHandshake  :", ti.TLSHandshake)
	fmt.Println("  ServerTime    :", ti.ServerTime)
	fmt.Println("  ResponseTime  :", ti.ResponseTime)
	fmt.Println("  TotalTime     :", ti.TotalTime)
	fmt.Println("  IsConnReused  :", ti.IsConnReused)
	fmt.Println("  IsConnWasIdle :", ti.IsConnWasIdle)
	fmt.Println("  ConnIdleTime  :", ti.ConnIdleTime)
	fmt.Println("  RequestAttempt:", ti.RequestAttempt)
	fmt.Println("  RemoteAddr    :", ti.RemoteAddr.String())
}
