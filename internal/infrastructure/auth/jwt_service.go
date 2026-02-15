package auth

import (
	"errors"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

type JWTService struct {
	secretKey string
	issuer    string
	duration  time.Duration
}

func NewJWTService(secret string) *JWTService {
	return &JWTService{
		secretKey: secret,
		issuer:    "workout-tracker",
		duration:  24 * time.Hour,
	}
}

func (j *JWTService) Generate(userID string) (string, time.Time, error) {
	if j == nil {
		return "", time.Time{}, errors.New("jwt service is nil")
	}
	if j.secretKey == "" {
		return "", time.Time{}, errors.New("jwt secret key is empty")
	}

	now := time.Now()
	exp := now.Add(j.duration)

	claims := jwt.RegisteredClaims{
		Subject:   userID,
		Issuer:    j.issuer,
		IssuedAt:  jwt.NewNumericDate(now),
		ExpiresAt: jwt.NewNumericDate(exp),
	}

	t := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	s, err := t.SignedString([]byte(j.secretKey))
	if err != nil {
		return "", time.Time{}, err
	}

	return s, exp, nil
}
