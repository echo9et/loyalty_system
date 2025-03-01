package interfaces

import "github.com/gin-gonic/gin"

func Midleware(ctx *gin.Context) {
	print("midleware 1")
	ctx.Next()
}

func Midleware2(ctx *gin.Context) {
	print("midleware 2")
	ctx.Next()
}
