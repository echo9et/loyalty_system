package user

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"gophermart.ru/internal/entities"
)

func Balance(group *gin.RouterGroup, mngr entities.WalletManagment) {

	group.GET("", func(ctx *gin.Context) {
		IDUser := ctx.Value("id_user").(int)

		wallet, err := mngr.Balance(IDUser)
		if err != nil {
			ctx.AbortWithError(http.StatusInternalServerError, err)
			return
		}

		ctx.JSON(200, wallet)
	})

	group.POST("/withdraw", func(ctx *gin.Context) {
		IDUser := ctx.Value("id_user").(int)
		withdraw := entities.Withdraw{ID: IDUser}
		if err := ctx.BindJSON(&withdraw); err != nil {
			ctx.AbortWithError(http.StatusBadRequest, err)
			return
		}
		err := mngr.Withdraw(withdraw)

		if err != nil {
			if err == entities.ErrNoMoney {
				ctx.AbortWithError(http.StatusPaymentRequired, err)
			} else if err == entities.ErrIncorrectOrder {
				ctx.AbortWithError(http.StatusUnprocessableEntity, err)
			} else {
				ctx.AbortWithError(http.StatusInternalServerError, err)
			}
			return
		}

		ctx.JSON(200, gin.H{
			"status": "ok",
		})
	})
}
