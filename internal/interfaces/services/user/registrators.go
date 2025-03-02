package user

import "github.com/gin-gonic/gin"

func Registrators(group *gin.RouterGroup) {
	group.POST("", func(ctx *gin.Context) {
		ctx.JSON(200, gin.H{
			"result": "ok /registrators",
		})
	})
}
