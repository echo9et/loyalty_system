package services

import (
	"github.com/gin-gonic/gin"
	user "gophermart.ru/internal/interfaces/services/user"
)

func User(group *gin.RouterGroup) {
	user.Login(group.Group("/login"))
	user.Registrators(group.Group("/registrators"))
	user.Balance(group.Group("/balance"))
	user.Orders(group.Group("/orders"))
	user.Withdrawals(group.Group("/withdrawals"))
}
