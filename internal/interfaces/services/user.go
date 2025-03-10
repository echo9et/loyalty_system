package services

import (
	"github.com/gin-gonic/gin"
	"gophermart.ru/internal/entities"
	user "gophermart.ru/internal/interfaces/services/user"
)

func User(group *gin.RouterGroup, managment entities.UserManagment) {
	user.Login(group.Group("/login"), managment)
	user.Register(group.Group("/register"), managment)
	user.Balance(group.Group("/balance"))
	user.Orders(group.Group("/orders"))
	user.Withdrawals(group.Group("/withdrawals"))
}
