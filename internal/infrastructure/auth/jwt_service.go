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

func (j *JWTService) Validate(tokenString string) (string, error) {
	if j == nil {
		return "", errors.New("jwt service is nil")
	}
	if j.secretKey == "" {
		return "", errors.New("jwt secret key is empty")
	}

	claims := &jwt.RegisteredClaims{}
	_, err := jwt.ParseWithClaims(
		tokenString,
		claims,
		func(t *jwt.Token) (interface{}, error) {
			if t.Method != jwt.SigningMethodHS256 {
				return nil, errors.New("unexpected signing method")
			}
			return []byte(j.secretKey), nil
		},
		jwt.WithValidMethods([]string{jwt.SigningMethodHS256.Alg()}),
		jwt.WithIssuer(j.issuer),
	)
	if err != nil {
		return "", err
	}

	if claims.Subject == "" {
		return "", errors.New("token subject is empty")
	}

	return claims.Subject, nil
}
