package jwt

import (
	"crypto/rsa"
	"os"
	"path/filepath"
	"time"

	"github.com/golang-jwt/jwt/v5"

	"loki/internal/app/errors"
	"loki/internal/config"
)

const (
	Dir            = "jwt"
	PrivateKeyFile = "private.key"
	PublicKeyFile  = "public.key"
)

type Payload struct {
	ID          string   `json:"id"`
	Roles       []string `json:"roles,omitempty"`
	Permissions []string `json:"permissions,omitempty"`
	Scope       []string `json:"scope,omitempty"`
}

type Jwt interface {
	Generate(payload Payload, duration time.Duration) (string, error)
	Verify(token string) (bool, error)
	Decode(token string) (*Payload, error)
}

type jwtService struct {
	cfg        *config.Config
	privateKey *rsa.PrivateKey
	publicKey  *rsa.PublicKey
}

type Claims struct {
	jwt.RegisteredClaims
	Roles       []string `json:"roles,omitempty"`
	Permissions []string `json:"permissions,omitempty"`
	Scope       []string `json:"scope,omitempty"`
}

func NewJWT(cfg *config.Config) (Jwt, error) {
	privateKey, publicKey, err := loadKeys(cfg)
	if err != nil {
		return nil, err
	}

	return &jwtService{
		cfg:        cfg,
		privateKey: privateKey,
		publicKey:  publicKey,
	}, nil
}

func (j *jwtService) Generate(payload Payload, duration time.Duration) (string, error) {
	claims := Claims{
		RegisteredClaims: jwt.RegisteredClaims{
			ID:        payload.ID,
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(duration)),
		},
		Roles:       payload.Roles,
		Permissions: payload.Permissions,
		Scope:       payload.Scope,
	}

	token := jwt.NewWithClaims(jwt.SigningMethodRS256, claims)

	signedToken, err := token.SignedString(j.privateKey)
	if err != nil {
		return "", err
	}

	return signedToken, nil
}

func (j *jwtService) Verify(token string) (bool, error) {
	claims := &Claims{}

	result, err := jwt.ParseWithClaims(token, claims,
		func(t *jwt.Token) (interface{}, error) {
			if _, ok := t.Method.(*jwt.SigningMethodRSA); !ok {
				return false, errors.ErrInvalidSigningMethod
			}
			return j.publicKey, nil
		})

	if err != nil {
		return false, err
	}

	if !result.Valid {
		return false, errors.ErrInvalidToken
	}

	return true, nil
}

func (j *jwtService) Decode(token string) (*Payload, error) {
	claims := &Claims{}

	result, err := jwt.ParseWithClaims(token, claims,
		func(t *jwt.Token) (interface{}, error) {
			if _, ok := t.Method.(*jwt.SigningMethodRSA); !ok {
				return false, errors.ErrInvalidSigningMethod
			}
			return j.publicKey, nil
		})

	if err != nil {
		return nil, err
	}

	if !result.Valid {
		return nil, errors.ErrInvalidToken
	}

	return &Payload{
		ID:          claims.ID,
		Roles:       claims.Roles,
		Permissions: claims.Permissions,
		Scope:       claims.Scope,
	}, nil
}

func loadKeys(cfg *config.Config) (*rsa.PrivateKey, *rsa.PublicKey, error) {
	privateKey, err := loadPrivateKey(cfg)
	if err != nil {
		return nil, nil, errors.ErrPrivateKeyNotFound
	}

	publicKey, err := loadPublicKey(cfg)
	if err != nil {
		return nil, nil, errors.ErrPublicKeyNotFound
	}

	return privateKey, publicKey, nil
}

func loadPrivateKey(cfg *config.Config) (*rsa.PrivateKey, error) {
	filePath := filepath.Join(cfg.CertPath, Dir, PrivateKeyFile)
	bytes, err := os.ReadFile(filePath)
	if err != nil {
		return nil, err
	}

	key, err := jwt.ParseRSAPrivateKeyFromPEM(bytes)
	if err != nil {
		return nil, err
	}

	return key, nil
}

func loadPublicKey(cfg *config.Config) (*rsa.PublicKey, error) {
	filePath := filepath.Join(cfg.CertPath, Dir, PublicKeyFile)
	bytes, err := os.ReadFile(filePath)
	if err != nil {
		return nil, err
	}

	key, err := jwt.ParseRSAPublicKeyFromPEM(bytes)
	if err != nil {
		return nil, err
	}

	return key, nil
}
