package services

import "github.com/gin-gonic/gin"

func Login(group *gin.RouterGroup) {
	group.POST("", func(ctx *gin.Context) {
		ctx.JSON(200, gin.H{
			"result": "ok /Login",
		})
	})
}
