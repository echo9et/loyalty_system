package user

import (
	"github.com/gin-gonic/gin"
)

func Balance(group *gin.RouterGroup) {
	group.GET("", func(ctx *gin.Context) {
		ctx.JSON(200, gin.H{
			"result": "ok",
		})
	})

	group.POST("/withdraw", func(ctx *gin.Context) {
		ctx.JSON(200, gin.H{
			"result": "ok /withdraw",
		})
	})
}
