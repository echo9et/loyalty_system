package utils

import (
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"strconv"
	"unicode"

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

	return claims.IDUser, nil
}

func IsValidOrder(number string) bool {
	for _, char := range number {
		if !unicode.IsDigit(char) {
			return false
		}
	}
	return isValidLuhn(number)
}

func isValidLuhn(number string) bool {
	sum := 0
	isEven := false
	for i := len(number) - 1; i >= 0; i-- {
		digit, _ := strconv.Atoi(string(number[i]))
		if isEven {
			digit *= 2
			if digit > 9 {
				digit -= 9
			}
		}
		sum += digit
		isEven = !isEven
	}
	return sum%10 == 0
}
