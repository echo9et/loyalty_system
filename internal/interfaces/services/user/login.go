package user

import (
	"errors"
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
			ctx.AbortWithError(http.StatusBadRequest, err)
			return
		}

		if user.Login == "" || user.HashPassword == "" {
			ctx.AbortWithError(http.StatusBadRequest,
				errors.New("неверный формат запроса"))
			return
		}
		user.HashPassword = utils.Sha256hash(user.HashPassword)

		u, err := mngr.User(user.Login)
		if err != nil {
			ctx.AbortWithError(http.StatusInternalServerError, err)
			return
		}

		if u == nil {
			ctx.AbortWithError(http.StatusUnauthorized,
				errors.New("неверная пара логин/пароль"))
			return
		}

		if !user.IsEcual(u) {
			ctx.AbortWithError(http.StatusUnauthorized,
				errors.New("неверная пара логин/пароль"))
			return
		}

		token, err := getToken(u)
		if err != nil {
			ctx.AbortWithError(http.StatusInternalServerError, err)
			return
		}
		ctx.SetCookie("token", token, int(config.Get().AliveToken), "/", "", false, true)
		ctx.JSON(http.StatusOK, gin.H{
			"status": "ok",
		})
	})
}
