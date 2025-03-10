package user

import (
	"time"

	"github.com/golang-jwt/jwt/v5"
	config "gophermart.ru/internal"
	"gophermart.ru/internal/entities"
)

func getToken(u *entities.User) (string, error) {
	expirationTime := time.Now().Add(config.Get().AliveToken)
	claims := &Claims{
		Login: u.Login,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(expirationTime),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString([]byte(config.Get().SecretKey))
	if err != nil {
		return "", err
	}

	return tokenString, nil
}
