package user

import (
	"fmt"
	"log/slog"
	"time"

	"github.com/golang-jwt/jwt/v5"
	config "gophermart.ru/internal"
	"gophermart.ru/internal/entities"
)

func getToken(u *entities.User) (string, error) {
	expirationTime := time.Now().Add(config.Get().AliveToken)
	claims := &entities.Claims{
		IDUser: u.ID,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(expirationTime),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString([]byte(config.Get().SecretKey))
	if err != nil {
		slog.Error(fmt.Sprintf("getToken %s", err))
		return "", err
	}

	return tokenString, nil
}
