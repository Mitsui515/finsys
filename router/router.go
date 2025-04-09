package router

import (
	"context"

	"github.com/Mitsui515/finsys/controller"
	"github.com/Mitsui515/finsys/middleware"
	"github.com/cloudwego/hertz/pkg/app"
	"github.com/cloudwego/hertz/pkg/app/server"
	"github.com/cloudwego/hertz/pkg/common/utils"
	"github.com/cloudwego/hertz/pkg/protocol/consts"
)

func RegisterRoutes(h *server.Hertz) {
	h.GET("/ping", func(ctx context.Context, c *app.RequestContext) {
		c.JSON(consts.StatusOK, utils.H{"message": "pong"})
	})
	transactionController := controller.NewTransactionController()
	userController := controller.NewUserController()
	fraudReportController := controller.NewFraudReportController()
	api := h.Group("/api")
	{
		auth := api.Group("/auth")
		{
			auth.POST("/register", userController.Register)
			auth.POST("/login", userController.Login)
		}
		user := api.Group("/user", middleware.JWTAuth())
		{
			user.GET("/info", userController.GetUserInfo)
		}
		transactions := api.Group("/transactions", middleware.JWTAuth())
		{
			transactions.GET("", transactionController.ListTransactionHandler)
			transactions.GET("/:id", transactionController.GetTransactionHandler)
			transactions.POST("", transactionController.CreateTransactionHandler)
			transactions.PUT("/:id", transactionController.UpdateTransactionHandler)
			transactions.DELETE("/:id", transactionController.DeleteTransactionHandler)
			transactions.POST("/import", transactionController.ImportTransactionsHandler)
		}
		fraudReports := api.Group("/fraud-reports", middleware.JWTAuth())
		{
			fraudReports.GET("", fraudReportController.ListFraudReportsHandler)
			fraudReports.GET("/:id", fraudReportController.GetFraudReportHandler)
			fraudReports.GET("/transaction/:transaction_id", fraudReportController.GetFraudReportByTransactionHandler)
			fraudReports.POST("", fraudReportController.CreateFraudReportHandler)
			fraudReports.PUT("/:id", fraudReportController.UpdateFraudReportHandler)
			fraudReports.DELETE("/:id", fraudReportController.DeleteFraudReportHandler)
			fraudReports.POST("/generate/:transaction_id", fraudReportController.GenerateReportHandler)
		}
	}
}
