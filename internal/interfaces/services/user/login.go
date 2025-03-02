package user

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	config "gophermart.ru/internal"
)

type Credentials struct {
	Password string `json:"password"`
	Username string `json:"username"`
}

type Claims struct {
	Username string `json:"username"`
	jwt.RegisteredClaims
}

func Login(group *gin.RouterGroup) {
	group.POST("", func(ctx *gin.Context) {
		var creds Credentials

		if err := ctx.BindJSON(&creds); err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{
				"error": "Bad Request",
			})
			return
		}

		expirationTime := time.Now().Add(config.Get().AliveToken)
		claims := &Claims{
			Username: creds.Username,
			RegisteredClaims: jwt.RegisteredClaims{
				ExpiresAt: jwt.NewNumericDate(expirationTime),
			},
		}

		token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
		tokenString, err := token.SignedString(config.Get().SecretKey)
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{
				"error": "Status Internal Server",
			})
			return
		}

		ctx.SetCookie("token", tokenString, int(config.Get().AliveToken), "", "", false, true)

		ctx.JSON(200, gin.H{
			"result": "success",
		})
	})
}
