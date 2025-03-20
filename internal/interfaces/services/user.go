package services

import (
	"github.com/gin-gonic/gin"
	"gophermart.ru/internal/infrastructure/storage"
	user "gophermart.ru/internal/interfaces/services/user"
)

func User(group *gin.RouterGroup, db *storage.Database) {
	user.Login(group.Group("/login"), db)
	user.Register(group.Group("/register"), db)
	user.Balance(group.Group("/balance"), db)
	user.Orders(group.Group("/orders"), db)
	user.Withdrawals(group.Group("/withdrawals"), db)
}
