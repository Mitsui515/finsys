package controller

import (
	"context"
	"path/filepath"
	"strconv"
	"time"

	"github.com/Mitsui515/finsys/config"
	"github.com/Mitsui515/finsys/service"
	"github.com/cloudwego/hertz/pkg/app"
	"github.com/cloudwego/hertz/pkg/common/utils"
	"github.com/cloudwego/hertz/pkg/protocol/consts"
)

type TransactionController struct {
	transactionService *service.TransactionService
}

func NewTransactionController() *TransactionController {
	return &TransactionController{
		transactionService: service.NewTransactionService(config.DB),
	}
}

func (c *TransactionController) GetTransactionHandler(ctx context.Context, reqCtx *app.RequestContext) {
	idStr := reqCtx.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		reqCtx.JSON(consts.StatusBadRequest, utils.H{
			"code":    consts.StatusBadRequest,
			"message": "Bad Request",
			"details": "Invalid transaction ID",
		})
		return
	}
	transaction, err := c.transactionService.GetByID(uint(id))
	if err != nil {
		reqCtx.JSON(consts.StatusNotFound, utils.H{
			"code":    consts.StatusNotFound,
			"message": "Not Found",
			"details": "Transaction not found",
		})
		return
	}
	reqCtx.JSON(consts.StatusOK, transaction)
}

func (c *TransactionController) ListTransactionHandler(ctx context.Context, reqCtx *app.RequestContext) {
	pageStr := reqCtx.Query("page")
	sizeStr := reqCtx.Query("size")
	transactionType := reqCtx.Query("type")
	startTimeStr := reqCtx.Query("start_time")
	endTimeStr := reqCtx.Query("end_time")
	page, err := strconv.Atoi(pageStr)
	if err != nil || page <= 0 {
		page = 1
	}
	size, err := strconv.Atoi(sizeStr)
	if err != nil || size <= 0 {
		size = 10
	}
	var startTime, endTime *time.Time
	if startTimeStr != "" {
		t, err := time.Parse(time.RFC3339, startTimeStr)
		if err != nil {
			reqCtx.JSON(consts.StatusBadRequest, utils.H{
				"code":    consts.StatusBadRequest,
				"message": "Bad Request",
				"details": "Invalid start time format",
			})
			return
		}
		startTime = &t
	}
	if endTimeStr != "" {
		t, err := time.Parse(time.RFC3339, endTimeStr)
		if err != nil {
			reqCtx.JSON(consts.StatusBadRequest, utils.H{
				"code":    consts.StatusBadRequest,
				"message": "Bad Request",
				"details": "Invalid end time format",
			})
			return
		}
		endTime = &t
	}
	transactions, err := c.transactionService.ListByPage(page, size, transactionType, startTime, endTime)
	if err != nil {
		reqCtx.JSON(consts.StatusInternalServerError, utils.H{
			"code":    consts.StatusInternalServerError,
			"message": "Internal Server Error",
			"details": err.Error(),
		})
		return
	}
	reqCtx.JSON(consts.StatusOK, transactions)
}

func (c *TransactionController) CreateTransactionHandler(ctx context.Context, reqCtx *app.RequestContext) {
	var req service.TransactionRequest
	if err := reqCtx.BindJSON(&req); err != nil {
		reqCtx.JSON(consts.StatusBadRequest, utils.H{
			"code":    consts.StatusBadRequest,
			"message": "Bad Request",
			"details": "Invalid request body",
		})
		return
	}
	id, err := c.transactionService.Create(&req)
	if err != nil {
		reqCtx.JSON(consts.StatusInternalServerError, utils.H{
			"code":    consts.StatusInternalServerError,
			"message": "Internal Server Error",
			"details": err.Error(),
		})
		return
	}
	reqCtx.JSON(consts.StatusOK, utils.H{
		"id": id,
	})
}

func (c *TransactionController) UpdateTransactionHandler(ctx context.Context, reqCtx *app.RequestContext) {
	idStr := reqCtx.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		reqCtx.JSON(consts.StatusBadRequest, utils.H{
			"code":    consts.StatusBadRequest,
			"message": "Bad Request",
			"details": "Invalid transaction ID",
		})
		return
	}
	var req service.TransactionRequest
	if err := reqCtx.BindAndValidate(&req); err != nil {
		reqCtx.JSON(consts.StatusBadRequest, utils.H{
			"code":    consts.StatusBadRequest,
			"message": "Bad Request",
			"details": "Invalid transaction data",
		})
		return
	}
	transaction, err := c.transactionService.Update(uint(id), &req)
	if err != nil {
		reqCtx.JSON(consts.StatusNotFound, utils.H{
			"code":    consts.StatusNotFound,
			"message": "Not Found",
			"details": "Transaction not found",
		})
		return
	}
	reqCtx.JSON(consts.StatusOK, transaction)
}

func (c *TransactionController) DeleteTransactionHandler(ctx context.Context, reqCtx *app.RequestContext) {
	idStr := reqCtx.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		reqCtx.JSON(consts.StatusBadRequest, utils.H{
			"code":    consts.StatusBadRequest,
			"message": "Bad Request",
			"details": "Invalid transaction ID",
		})
		return
	}
	err = c.transactionService.Delete(uint(id))
	if err != nil {
		reqCtx.JSON(consts.StatusNotFound, utils.H{
			"code":    consts.StatusNotFound,
			"message": "Not Found",
			"details": "Transaction not found",
		})
		return
	}
	reqCtx.Status(consts.StatusNoContent)
}

func (c *TransactionController) ImportTransactionsHandler(ctx context.Context, reqCtx *app.RequestContext) {
	file, err := reqCtx.FormFile("file")
	if err != nil {
		reqCtx.JSON(consts.StatusBadRequest, utils.H{
			"code":    consts.StatusBadRequest,
			"message": "Bad Request",
			"details": "Please select CSV file to upload",
		})
		return
	}
	if filepath.Ext(file.Filename) != ".csv" {
		reqCtx.JSON(consts.StatusBadRequest, utils.H{
			"code":    consts.StatusBadRequest,
			"message": "Bad Request",
			"details": "Only CSV file is allowed",
		})
		return
	}
	src, err := file.Open()
	if err != nil {
		reqCtx.JSON(consts.StatusInternalServerError, utils.H{
			"code":    consts.StatusInternalServerError,
			"message": "Internal Server Error",
			"details": "Cannot read file",
		})
		return
	}
	defer src.Close()
	importService := service.NewTransactionService(config.DB)
	count, err := importService.ImportFromCSV(src)
	if err != nil {
		reqCtx.JSON(consts.StatusInternalServerError, utils.H{
			"code":    consts.StatusInternalServerError,
			"message": "Internal Server Error",
			"details": err.Error(),
		})
		return
	}
	reqCtx.JSON(consts.StatusOK, utils.H{
		"message": "Successfully imported transactions",
		"count":   count,
	})
}
