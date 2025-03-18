package user

import (
	"log/slog"
	"net/http"

	"github.com/gin-gonic/gin"
	config "gophermart.ru/internal"
	"gophermart.ru/internal/entities"
	"gophermart.ru/internal/utils"
)

func Register(group *gin.RouterGroup, mngr entities.UserManagment) {
	group.POST("", func(ctx *gin.Context) {
		var user entities.User

		if err := ctx.BindJSON(&user); err != nil {
			slog.Error(err.Error())
			ctx.JSON(http.StatusBadRequest, gin.H{
				"error": "неверный формат запроса",
			})
			return
		}

		if user.Login == "" || user.HashPassword == "" {
			slog.Error("Переданы пустые поля")
			ctx.JSON(http.StatusBadRequest, gin.H{
				"error": "неверный формат запроса",
			})
			return
		}
		user.HashPassword = utils.Sha256hash(user.HashPassword)

		if err := mngr.InsertUser(user); err != nil {
			slog.Error(err.Error())
			ctx.JSON(http.StatusConflict, gin.H{
				"error": "логин уже занят",
			})
			return
		}

		token, err := getToken(&user)
		if err != nil {
			slog.Error(err.Error())
			ctx.JSON(http.StatusInternalServerError, gin.H{
				"error": "внутренняя ошибка сервера",
			})
			return
		}

		ctx.SetCookie("token", token, int(config.Get().AliveToken), "", "", false, true)

		ctx.JSON(200, gin.H{
			"result": "ok",
		})
	})
}
