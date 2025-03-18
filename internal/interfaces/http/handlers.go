package interfaces

import (
	"errors"
	"log/slog"
	"net/http"

	"github.com/gin-gonic/gin"
)

func MidlewareAuth(ctx *gin.Context) {
	if ctx.Request.URL.Path == "/api/user/login" || ctx.Request.URL.Path == "/api/user/register" {
		ctx.Next()
		return
	}

	_, err := ctx.Cookie("token")
	if err != nil {
		slog.Error(err.Error())
		ctx.AbortWithError(http.StatusUnauthorized, errors.New("no unauthorized"))
		return
	}

	ctx.Next()
}

func Midleware2(ctx *gin.Context) {
	ctx.Next()
}
