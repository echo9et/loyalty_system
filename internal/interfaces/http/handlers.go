package interfaces

import (
	"errors"
	"fmt"
	"log/slog"
	"net/http"

	"github.com/gin-gonic/gin"
	"gophermart.ru/internal/utils"
)

func MidlewareAuth(ctx *gin.Context) {
	if ctx.Request.URL.Path == "/api/user/login" || ctx.Request.URL.Path == "/api/user/register" {
		ctx.Next()
		return
	}

	token, err := ctx.Cookie("token")
	if err != nil {
		slog.Error(fmt.Sprintf("MidlewareAuth: %s", err.Error()))
		ctx.AbortWithError(http.StatusUnauthorized, errors.New("no unauthorized"))
		return
	}

	IDUser, err := utils.LoginFromToken(token)

	if err != nil {
		slog.Error(err.Error())
		ctx.AbortWithError(http.StatusUnauthorized, errors.New("no unauthorized"))
	}

	ctx.Set("id_user", IDUser)
	fmt.Println("id_user", IDUser)
	ctx.Next()
}

func Midleware2(ctx *gin.Context) {
	ctx.Next()
}
