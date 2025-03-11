package utils

import (
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"fmt"

	"github.com/golang-jwt/jwt/v5"
	config "gophermart.ru/internal"
	"gophermart.ru/internal/entities"
)

func Sha256hash(value string) string {
	hash := sha256.New()
	hash.Write([]byte(value))
	return hex.EncodeToString(hash.Sum(nil))
}

func LoginFromToken(sToken string) (int, error) {
	claims := &entities.Claims{}
	tkn, err := jwt.ParseWithClaims(sToken, claims, func(jwtKey *jwt.Token) (any, error) {
		return []byte(config.Get().SecretKey), nil
	})
	if err != nil {
		return -1, err
	}

	if !tkn.Valid {
		return -1, errors.New("no valid token")
	}
	fmt.Println("OK")
	return claims.IdUser, nil
}
