package user

import (
	"fmt"
	"log/slog"
	"net/http"

	"github.com/gin-gonic/gin"
	"gophermart.ru/internal/entities"
)

func Withdrawals(group *gin.RouterGroup, mngr entities.WalletManagment) {
	group.GET("", func(ctx *gin.Context) {
		IDUser := ctx.Value("id_user").(int)

		withdraws, err := mngr.Withdraws(IDUser)

		if err != nil {
			slog.Error(fmt.Sprintf("GET Withdrawals %s", err.Error()))
			ctx.AbortWithStatus(http.StatusInternalServerError)
			return
		}

		if len(withdraws) == 0 {
			ctx.AbortWithStatus(http.StatusNoContent)
			return
		}

		ctx.JSON(200, withdraws)
	})
}
