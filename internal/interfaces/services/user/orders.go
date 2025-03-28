package user

import (
	"errors"
	"fmt"
	"io"
	"log/slog"
	"net/http"

	"github.com/gin-gonic/gin"
	"gophermart.ru/internal/entities"
	"gophermart.ru/internal/utils"
)

func Orders(group *gin.RouterGroup, mngr entities.OrdersManagment, a *AccrualSystem) {
	group.POST("", func(ctx *gin.Context) {
		body, err := io.ReadAll(ctx.Request.Body)
		if err != nil {
			slog.Error(err.Error())
			ctx.JSON(http.StatusBadRequest, gin.H{})
			return
		}

		number := string(body)
		if !utils.IsValidOrder(number) {
			slog.Error("ошибка валидации номера ордера")
			ctx.JSON(http.StatusUnprocessableEntity, gin.H{})
			return
		}

		order, err := mngr.Order(number)

		if err != nil {
			slog.Error(fmt.Sprintf("ошибка получения ордера %s", err))
			ctx.JSON(http.StatusBadRequest, gin.H{})
			return
		}

		IDUser := ctx.Value("id_user").(int)

		if order != nil {
			if order.IDUser == IDUser {
				ctx.JSON(http.StatusOK, gin.H{
					"answer": "номер заказа уже был загружен этим пользователем"})
			} else {
				ctx.AbortWithError(http.StatusConflict,
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
		IDUser := ctx.Value("id_user").(int)

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
