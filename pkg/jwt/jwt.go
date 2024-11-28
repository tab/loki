package jwt

import (
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"

	"loki/internal/app/errors"
	"loki/internal/config"
)

const TokenExp = time.Minute * 30

type Jwt interface {
	Generate(id uuid.UUID) (string, error)
	Verify(token string) (bool, error)
}

type jwtService struct {
	cfg *config.Config
}

type Claims struct {
	jwt.RegisteredClaims
	UserID uuid.UUID
}

func NewJWTService(cfg *config.Config) Jwt {
	return &jwtService{cfg: cfg}
}

func (j *jwtService) Generate(id uuid.UUID) (string, error) {
	result := jwt.NewWithClaims(jwt.SigningMethodHS256, Claims{
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(TokenExp)),
		},
		UserID: id,
	})

	token, err := result.SignedString([]byte(j.cfg.SecretKey))
	if err != nil {
		return "", err
	}

	return token, nil
}

func (j *jwtService) Verify(token string) (bool, error) {
	claims := &Claims{}

	result, err := jwt.ParseWithClaims(token, claims,
		func(t *jwt.Token) (interface{}, error) {
			if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
				return false, errors.ErrInvalidSigningMethod
			}
			return []byte(j.cfg.SecretKey), nil
		})

	if err != nil {
		return false, err
	}

	if !result.Valid {
		return false, errors.ErrInvalidToken
	}

	return true, nil
}
