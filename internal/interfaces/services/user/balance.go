package user

import (
	"fmt"
	"log/slog"
	"net/http"

	"github.com/gin-gonic/gin"
	"gophermart.ru/internal/entities"
)

func Balance(group *gin.RouterGroup, mngr entities.WalletManagment) {

	group.GET("", func(ctx *gin.Context) {
		IDUser := ctx.Value("id_user").(int)

		wallet, err := mngr.Balance(IDUser)
		if err != nil {
			slog.Error(fmt.Sprintf("GET Balance %s", err.Error()))
			ctx.AbortWithStatus(http.StatusInternalServerError)
			return
		}

		ctx.JSON(200, wallet)
	})

	group.POST("/withdraw", func(ctx *gin.Context) {
		IDUser := ctx.Value("id_user").(int)
		withdraw := entities.Withdraw{ID: IDUser}
		if err := ctx.BindJSON(&withdraw); err != nil {
			slog.Error(fmt.Sprintf("POST Balance %s", err.Error()))
			ctx.AbortWithStatus(http.StatusBadRequest)
		}
		err := mngr.Withdraw(withdraw)

		if err != nil {
			slog.Error(fmt.Sprintf("POST Balance %s", err.Error()))
			if err == entities.ErrNoMoney {
				ctx.AbortWithStatus(http.StatusPaymentRequired)
			} else if err == entities.ErrIncorrectOrder {
				ctx.AbortWithStatus(http.StatusUnprocessableEntity)
			} else {
				ctx.AbortWithStatus(http.StatusInternalServerError)
			}
			return
		}

		ctx.JSON(200, gin.H{
			"status": "ok",
		})
	})
}
