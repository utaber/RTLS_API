package auth

import (
	"time"

	"github.com/golang-jwt/jwt/v5"
)

type JWTService struct {
	Secret string
	Algo   jwt.SigningMethod
	Expire time.Duration
}

func NewJWTService(secret string, algo jwt.SigningMethod, expire time.Duration) *JWTService {
	return &JWTService{secret, algo, expire}
}

func (j *JWTService) GenerateToken() (string, error) {
	claims := jwt.MapClaims{
		"sub": "admin",
		"exp": time.Now().Add(j.Expire).Unix(),
	}

	token := jwt.NewWithClaims(j.Algo, claims)
	return token.SignedString([]byte(j.Secret))
}
