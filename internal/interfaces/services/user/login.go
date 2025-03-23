package user

import (
	"log/slog"
	"net/http"

	"github.com/gin-gonic/gin"
	config "gophermart.ru/internal"
	"gophermart.ru/internal/entities"
	"gophermart.ru/internal/utils"
)

func Login(group *gin.RouterGroup, mngr entities.UserManagment) {
	group.POST("", func(ctx *gin.Context) {
		var user entities.User

		if err := ctx.BindJSON(&user); err != nil {
			slog.Error(err.Error())
			ctx.SetCookie("token", "", 0, "/", "", false, true)
			ctx.JSON(http.StatusBadRequest, gin.H{
				"error": "Bad Request",
			})
			return
		}

		if user.Login == "" || user.HashPassword == "" {
			ctx.SetCookie("token", "", 0, "/", "", false, true)
			slog.Error("Переданы пустые поля")
			ctx.JSON(http.StatusBadRequest, gin.H{
				"error": "неверный формат запроса",
			})
			return
		}
		user.HashPassword = utils.Sha256hash(user.HashPassword)

		u, err := mngr.User(user.Login)
		if err != nil {
			ctx.SetCookie("token", "", 0, "/", "", false, true)
			ctx.JSON(http.StatusInternalServerError, gin.H{
				"error": "внутренняя ошибка сервера",
			})
			return
		}

		if u == nil {
			ctx.SetCookie("token", "", 0, "/", "", false, true)
			ctx.JSON(http.StatusUnauthorized, gin.H{
				"error": "неверная пара логин/пароль",
			})
			return
		}

		if !user.IsEcual(u) {
			ctx.SetCookie("token", "", 0, "/", "", false, true)
			ctx.JSON(http.StatusUnauthorized, gin.H{
				"error": "неверная пара логин/пароль",
			})
			return
		}

		token, err := getToken(u)
		if err != nil {
			ctx.SetCookie("token", "", 0, "/", "", false, true)
			slog.Error(err.Error())
			ctx.JSON(http.StatusInternalServerError, gin.H{
				"error": "внутренняя ошибка сервера",
			})
			return
		}
		ctx.SetCookie("token", token, int(config.Get().AliveToken), "/", "", false, true)
		ctx.JSON(200, gin.H{
			"status": "ok",
		})
	})
}
