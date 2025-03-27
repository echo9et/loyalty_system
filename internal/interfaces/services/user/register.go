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
			ctx.AbortWithError(http.StatusConflict, err)
			return
		}

		u, err := mngr.User(user.Login)

		if err != nil || u == nil {
			ctx.AbortWithError(http.StatusInternalServerError, err)
			return
		}

		token, err := getToken(u)
		if err != nil {
			ctx.AbortWithError(http.StatusInternalServerError, err)
			return
		}

		ctx.SetCookie("token", token, int(config.Get().AliveToken), "", "", false, true)
		ctx.JSON(200, gin.H{"status": "ok"})
	})
}
