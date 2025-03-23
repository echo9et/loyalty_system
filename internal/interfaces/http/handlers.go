package interfaces

import (
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
		ctx.AbortWithError(http.StatusUnauthorized, err)
		return
	}

	IDUser, err := utils.LoginFromToken(token)

	if err != nil {
		ctx.AbortWithError(http.StatusUnauthorized, err)
		return
	}

	ctx.Set("id_user", IDUser)
	ctx.Next()
}

func MidlewareErrors(ctx *gin.Context) {
	ctx.Next()

	if len(ctx.Errors) > 0 {
		err := ctx.Errors[0]
		code := ctx.Writer.Status()

		if code == http.StatusInternalServerError {
			slog.Error(err.Err.Error())
		} else {
			slog.Warn(err.Err.Error())
		}
	}
}
