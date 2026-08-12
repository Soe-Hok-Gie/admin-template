package utils

import (
	"errors"
	"os"
	"strconv"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

func GenerateToken(ID int64) (string, error) {
	secret := []byte(os.Getenv("JWT_SECRET"))
	if len(secret) == 0 {
		return "", errors.New("JWT_SECRET_EMPTY")
	}

	now := time.Now()
	claims := jwt.RegisteredClaims{
		Issuer:    "service Login",
		Subject:   strconv.FormatInt(ID, 10),
		IssuedAt:  jwt.NewNumericDate(now),
		ExpiresAt: jwt.NewNumericDate(now.Add(15 * time.Minute)),
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(secret)
}

func ValidateToken(tokenstr string) (*int64, error) {
	secret := os.Getenv("JWT_SECRET")
	if secret == "" {
		return nil, errors.New("JWT_SECRET_EMPTY")
	}
	token, err := jwt.ParseWithClaims(tokenstr, &jwt.RegisteredClaims{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("expected signing method")
		}
		return []byte(secret), nil

	})
	if err != nil || !token.Valid {
		return nil, err
	}
	claims, ok := token.Claims.(*jwt.RegisteredClaims)
	if !ok || !token.Valid {
		return nil, errors.New("invalid token")
	}

	ID, err := strconv.ParseInt(claims.Subject, 10, 64)
	if err != nil {
		return nil, errors.New("invalid return error")
	}
	return &ID, nil
}
