package services

import "github.com/gin-gonic/gin"

func Withdrawals(group *gin.RouterGroup) {
	group.GET("", func(ctx *gin.Context) {
		ctx.JSON(200, gin.H{
			"result": "ok /Withdrawals",
		})
	})
}
