package interfaces

import (
	"log/slog"

	"github.com/gin-gonic/gin"
)

func Midleware(ctx *gin.Context) {
	slog.Info(ctx.Request.Method)
	ctx.Next()
}

func Midleware2(ctx *gin.Context) {
	print("midleware 2\n")
	ctx.Next()
}
