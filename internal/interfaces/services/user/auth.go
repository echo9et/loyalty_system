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
		IdUser: u.Id,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(expirationTime),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString([]byte(config.Get().SecretKey))
	if err != nil {
		slog.Error("getToken", err)
		return "", err
	}
	fmt.Println("getToken USER id: %d", claims.IdUser)

	return tokenString, nil
}
