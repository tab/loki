package services

import (
	"context"
	"crypto/x509"
	"encoding/base64"
	"fmt"
	"regexp"
	"strings"
	"time"

	"github.com/google/uuid"

	"loki/internal/app/errors"
	"loki/internal/app/models"
	"loki/internal/app/models/dto"
	"loki/internal/app/repositories"
	"loki/internal/app/serializers"
	"loki/internal/config"
	"loki/pkg/jwt"
	"loki/pkg/logger"
)

type Authentication interface {
	CreateSmartIdSession(ctx context.Context, params dto.CreateSmartIdSessionRequest) (*serializers.SmartIdSessionSerializer, error)
	GetSmartIdSessionStatus(ctx context.Context, id uuid.UUID) (*dto.SmartIdProviderSessionStatusResponse, error)
	UpdateSession(ctx context.Context, params dto.UpdateSmartIdSessionParams) (*serializers.SessionSerializer, error)
	FindSessionByID(ctx context.Context, id string) (*serializers.SessionSerializer, error)
	Complete(ctx context.Context, id string) (*serializers.UserSerializer, error)
}

type authentication struct {
	cfg             *config.Config
	smartIdProvider SmartIdProvider
	smartIDQueue    chan<- *SmartIdQueue
	database        repositories.Database
	redis           repositories.Redis
	jwt             jwt.Jwt
	log             *logger.Logger
}

func NewAuthentication(
	cfg *config.Config,
	smartIdProvider SmartIdProvider,
	smartIDQueue chan *SmartIdQueue,
	database repositories.Database,
	redis repositories.Redis,
	jwt jwt.Jwt,
	log *logger.Logger,
) Authentication {
	return &authentication{
		cfg:             cfg,
		smartIdProvider: smartIdProvider,
		smartIDQueue:    smartIDQueue,
		database:        database,
		redis:           redis,
		jwt:             jwt,
		log:             log,
	}
}

func (a *authentication) CreateSmartIdSession(ctx context.Context, params dto.CreateSmartIdSessionRequest) (response *serializers.SmartIdSessionSerializer, error error) {
	result, err := a.smartIdProvider.CreateSession(ctx, params)
	if err != nil {
		a.log.Error().Err(err).Msg("Failed to initiate SmartId authentication")
		return nil, err
	}

	id, err := uuid.Parse(result.ID)
	if err != nil {
		a.log.Error().Err(err).Msg("Invalid session ID format")
		return nil, err
	}

	err = a.redis.CreateSession(ctx, &models.Session{
		ID:           id,
		PersonalCode: params.PersonalCode,
		Code:         result.Code,
		Status:       models.SESSION_RUNNING,
	})
	if err != nil {
		a.log.Error().Err(err).Msg("Failed to create session")
		return nil, err
	}

	a.smartIDQueue <- &SmartIdQueue{
		ID: id,
	}

	return &serializers.SmartIdSessionSerializer{
		ID:   id,
		Code: result.Code,
	}, nil
}

func (a *authentication) GetSmartIdSessionStatus(_ context.Context, id uuid.UUID) (response *dto.SmartIdProviderSessionStatusResponse, error error) {
	result, err := a.smartIdProvider.GetSessionStatus(id)
	if err != nil {
		a.log.Error().Err(err).Msg("Failed to get session status")
		return nil, err
	}

	return result, nil
}

func (a *authentication) UpdateSession(ctx context.Context, params dto.UpdateSmartIdSessionParams) (response *serializers.SessionSerializer, error error) {
	err := a.redis.UpdateSession(ctx, &models.Session{
		ID:     params.ID,
		Status: params.Status,
		Payload: models.SessionPayload{
			State:     params.Payload.State,
			Result:    params.Payload.Result,
			Signature: params.Payload.Signature,
			Cert:      params.Payload.Cert,
		},
	})
	if err != nil {
		a.log.Error().Err(err).Msg("Failed to update session")
		return nil, err
	}

	return &serializers.SessionSerializer{
		ID:     params.ID,
		Status: params.Status,
	}, nil
}

func (a *authentication) FindSessionByID(ctx context.Context, sessionId string) (response *serializers.SessionSerializer, error error) {
	id, err := uuid.Parse(sessionId)
	if err != nil {
		a.log.Error().Err(err).Msg("Invalid session ID format")
		return nil, err
	}

	result, err := a.redis.FindSessionByID(ctx, id)
	if err != nil {
		a.log.Error().Err(err).Msg("Failed to find session")
		return nil, err
	}

	return &serializers.SessionSerializer{
		ID:     result.ID,
		Status: result.Status,
	}, nil
}

func (a *authentication) Complete(ctx context.Context, sessionId string) (response *serializers.UserSerializer, error error) {
	id, err := uuid.Parse(sessionId)
	if err != nil {
		a.log.Error().Err(err).Msg("Invalid session ID format")
		return nil, err
	}

	result, err := a.redis.FindSessionByID(ctx, id)
	if err != nil {
		a.log.Error().Err(err).Msg("Failed to find session")
		return nil, err
	}

	user, err := a.createUser(ctx, result.Payload)
	if err != nil {
		a.log.Error().Err(err).Msg("Failed to create user")
		return nil, err
	}

	err = a.redis.DeleteSessionByID(ctx, id)
	if err != nil {
		a.log.Error().Err(err).Msg("Failed to delete session")
		return nil, err
	}

	return &serializers.UserSerializer{
		ID:             user.ID,
		FirstName:      user.FirstName,
		LastName:       user.LastName,
		IdentityNumber: user.IdentityNumber,
		PersonalCode:   user.PersonalCode,
		AccessToken:    user.AccessToken,
		RefreshToken:   user.RefreshToken,
	}, nil
}

func (a *authentication) createUser(ctx context.Context, payload models.SessionPayload) (*models.User, error) {
	certificate, err := extractUserFromCertificate(payload.Cert.Value)
	if err != nil {
		a.log.Error().Err(err).Msg("Failed to extract user from certificate")
		return nil, err
	}

	accessToken, err := a.jwt.Generate(jwt.Payload{
		ID: certificate.IdentityNumber,
	}, models.AccessTokenExp)
	if err != nil {
		a.log.Error().Err(err).Msg("Failed to create access token")
		return nil, err
	}

	refreshToken, err := a.jwt.Generate(jwt.Payload{
		ID: certificate.IdentityNumber,
	}, models.RefreshTokenExp)
	if err != nil {
		a.log.Error().Err(err).Msg("Failed to create refresh token")
		return nil, err
	}

	user, err := a.database.CreateOrUpdateUserWithTokens(ctx, dto.CreateUserParams{
		IdentityNumber: certificate.IdentityNumber,
		PersonalCode:   certificate.PersonalCode,
		FirstName:      certificate.FirstName,
		LastName:       certificate.LastName,
		AccessToken: dto.CreateTokenParams{
			Type:      models.AccessTokenType,
			Value:     accessToken,
			ExpiresAt: time.Now().Add(models.AccessTokenExp),
		},
		RefreshToken: dto.CreateTokenParams{
			Type:      models.RefreshTokenType,
			Value:     refreshToken,
			ExpiresAt: time.Now().Add(models.RefreshTokenExp),
		},
	})
	if err != nil {
		a.log.Error().Err(err).Msg("Failed to create user")
		return nil, err
	}

	return user, nil
}

func extractUserFromCertificate(certValue string) (*dto.SmartIdProviderCertificateExtract, error) {
	certBytes, err := base64.StdEncoding.DecodeString(certValue)
	if err != nil {
		return nil, err
	}

	cert, err := x509.ParseCertificate(certBytes)
	if err != nil {
		return nil, err
	}

	subject := cert.Subject

	commonName := subject.CommonName
	parts := strings.Split(commonName, ",")
	if len(parts) < 2 {
		return nil, fmt.Errorf("invalid CommonName format: %s", commonName)
	}

	personalCode, _ := extractPersonalCode(subject.SerialNumber)
	firstName := strings.TrimSpace(parts[0])
	lastName := strings.TrimSpace(parts[1])

	return &dto.SmartIdProviderCertificateExtract{
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
