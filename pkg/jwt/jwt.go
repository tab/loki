package jwt

import (
	"time"

	"github.com/golang-jwt/jwt/v5"

	"loki/internal/app/errors"
	"loki/internal/config"
)

type Payload struct {
	ID string `json:"id"`
}

type Jwt interface {
	Generate(payload Payload, duration time.Duration) (string, error)
	Verify(token string) (bool, error)
}

type jwtService struct {
	cfg *config.Config
}

type Claims struct {
	jwt.RegisteredClaims
}

func NewJWT(cfg *config.Config) Jwt {
	return &jwtService{cfg: cfg}
}

func (j *jwtService) Generate(payload Payload, duration time.Duration) (string, error) {
	claims := Claims{
		RegisteredClaims: jwt.RegisteredClaims{
			ID:        payload.ID,
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(duration)),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	signedToken, err := token.SignedString([]byte(j.cfg.SecretKey))
	if err != nil {
		return "", err
	}

	return signedToken, nil
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
