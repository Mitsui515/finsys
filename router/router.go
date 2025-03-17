package router

import (
	"context"

	"github.com/Mitsui515/finsys/controller"
	"github.com/cloudwego/hertz/pkg/app"
	"github.com/cloudwego/hertz/pkg/app/server"
	"github.com/cloudwego/hertz/pkg/common/utils"
	"github.com/cloudwego/hertz/pkg/protocol/consts"
)

func RegisterRoutes(h *server.Hertz) {
	h.GET("/ping", func(ctx context.Context, c *app.RequestContext) {
		c.JSON(consts.StatusOK, utils.H{"message": "pong"})
	})

	TransactionController := controller.NewTransactionController()

	api := h.Group("/api")
	{
		transactions := api.Group("/transactions")
		{
			transactions.GET("", TransactionController.ListTransactionHandler)
			transactions.GET("/:id", TransactionController.GetTransactionHandler)
			transactions.POST("", TransactionController.CreateTransactionHandler)
			transactions.PUT("/:id", TransactionController.UpdateTransactionHandler)
			transactions.DELETE("/:id", TransactionController.DeleteTransactionHandler)
			transactions.POST("/import", TransactionController.ImportTransactionsHandler)
		}
	}
}
