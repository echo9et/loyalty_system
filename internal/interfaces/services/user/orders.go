package user

import (
	"errors"
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"unicode"

	"github.com/gin-gonic/gin"
	"gophermart.ru/internal/entities"
	"gophermart.ru/internal/utils"
)

func isValid(number string) bool {
	for _, char := range number {
		if !unicode.IsDigit(char) {
			return false
		}
	}
	return true
}

func Orders(group *gin.RouterGroup, mngr entities.OrdersManagment) {
	group.POST("", func(ctx *gin.Context) {
		body, err := io.ReadAll(ctx.Request.Body)
		if err != nil {
			slog.Error(err.Error())
			ctx.JSON(http.StatusBadRequest, gin.H{})
			return
		}

		number := string(body)
		if !isValid(number) {
			slog.Error("ошибка валидации номера ордера")
			ctx.JSON(http.StatusBadRequest, gin.H{})
			return
		}

		order, err := mngr.Order(number)

		if err != nil {
			slog.Error(fmt.Sprintf("ошибка получения ордера %s", err))
			ctx.JSON(http.StatusBadRequest, gin.H{})
			return
		}

		token, err := ctx.Cookie("token")
		if err != nil {
			slog.Error(fmt.Sprintf("--- MidlewareAuth token %s", err))
			ctx.AbortWithError(http.StatusUnauthorized,
				errors.New("no unauthorized"))
			return
		}

		idUser, err := utils.LoginFromToken(token)
		if (err != nil) || (idUser == -1) {
			ctx.AbortWithError(http.StatusUnauthorized,
				errors.New("no unauthorized"))
			return
		}
		fmt.Println("ORDERS USER_ID", idUser)

		if order != nil {
			if order.IdUser == idUser {
				ctx.JSON(200, gin.H{
					"answer": "номер заказа уже был загружен этим пользователем"})
			} else {
				ctx.AbortWithError(http.StatusConflict,
					errors.New("номер заказа уже был загружен другим пользователем"))
			}
			return
		}

		newOrder := entities.Order{
			Number: number,
			IdUser: idUser,
			Status: "NEW",
		}

		err = mngr.AddOrder(newOrder)
		if err != nil {
			slog.Error(fmt.Sprintf("add order %s", err.Error()))
			ctx.AbortWithStatus(http.StatusInternalServerError)
			return
		}

		ctx.JSON(200, gin.H{
			"result": "ok",
		})
	})
}
