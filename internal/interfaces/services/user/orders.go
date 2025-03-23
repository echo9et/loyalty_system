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

func Orders(group *gin.RouterGroup, mngr entities.OrdersManagment, a *AccrualSystem) {
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

		IDUser, err := utils.LoginFromToken(token)
		if (err != nil) || (IDUser == -1) {
			ctx.AbortWithError(http.StatusUnauthorized,
				errors.New("no unauthorized"))
			return
		}

		if order != nil {
			if order.IDUser == IDUser {
				ctx.JSON(200, gin.H{
					"answer": "номер заказа уже был загружен этим пользователем"})
			} else {
				ctx.AbortWithError(http.StatusOK,
					errors.New("номер заказа уже был загружен другим пользователем"))
			}
			return
		}

		newOrder := entities.Order{
			Number: number,
			IDUser: IDUser,
			Status: entities.OrderNew,
		}

		err = mngr.AddOrder(newOrder)
		if err != nil {
			slog.Error(fmt.Sprintf("add order %s", err.Error()))
			ctx.AbortWithStatus(http.StatusInternalServerError)
			return
		}

		go UpdateOrder(&newOrder, mngr, a)

		ctx.JSON(http.StatusAccepted, gin.H{
			"result": "ok",
		})
	})

	group.GET("", func(ctx *gin.Context) {
		token, err := ctx.Cookie("token")
		if err != nil {
			slog.Error(fmt.Sprintf("MidlewareAuth token %s", err))
			ctx.AbortWithError(http.StatusUnauthorized,
				errors.New("пользователь не авторизован"))
			return
		}

		IDUser, err := utils.LoginFromToken(token)
		if (err != nil) || (IDUser == -1) {
			ctx.AbortWithError(http.StatusUnauthorized,
				errors.New("пользователь не авторизован"))
			return
		}

		orders, err := mngr.Orders(IDUser)

		if err != nil {
			slog.Error(fmt.Sprintf("Orders %s", err))
			ctx.AbortWithError(http.StatusInternalServerError,
				errors.New("внутренняя ошибка сервера"))
			return
		}

		if len(orders) == 0 {
			ctx.AbortWithStatus(http.StatusNoContent)
			return
		}

		for _, order := range orders {
			UpdateOrder(&order, mngr, a)
		}

		ctx.JSON(http.StatusOK, orders)
	})
}

func UpdateOrder(order *entities.Order, mngr entities.OrdersManagment, a *AccrualSystem) {
	if order.Status == entities.OrderProcessed || order.Status == entities.OrderInvalid {
		return
	}

	newOrder, status, err := a.GetOrderInfo(order.Number)

	if err != nil {
		slog.Error(fmt.Sprintf("GetOrderInfo error %s", err.Error()))
		return
	}

	if status != http.StatusOK {
		slog.Error(fmt.Sprintf("GetOrderInfo status %d", status))
		return
	}
	if newOrder.Status == order.Status {
		return
	}

	order.Status = newOrder.Status
	order.Accrual = newOrder.Accrual
	err = mngr.UpdateOrder(*order)

	if err != nil {
		slog.Error(fmt.Sprintf("update order %s", err.Error()))
	}
}
