package controller

import (
	"context"
	"strconv"

	"github.com/Mitsui515/finsys/config"
	"github.com/Mitsui515/finsys/service"
	"github.com/cloudwego/hertz/pkg/app"
	"github.com/cloudwego/hertz/pkg/common/utils"
	"github.com/cloudwego/hertz/pkg/protocol/consts"
)

type FraudReportController struct {
	fraudReportService *service.FraudReportService
}

func NewFraudReportController() *FraudReportController {
	return &FraudReportController{
		fraudReportService: service.NewFraudReportService(config.DB),
	}
}

func (c *FraudReportController) GetFraudReportHandler(ctx context.Context, reqCtx *app.RequestContext) {
	idStr := reqCtx.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		reqCtx.JSON(consts.StatusBadRequest, utils.H{
			"code":    consts.StatusBadRequest,
			"message": "Bad Request",
			"details": "Invalid fraud report ID",
		})
		return
	}
	report, err := c.fraudReportService.GetByID(uint(id))
	if err != nil {
		reqCtx.JSON(consts.StatusNotFound, utils.H{
			"code":    consts.StatusNotFound,
			"message": "Not Found",
			"details": "Fraud report not found",
		})
		return
	}
	reqCtx.JSON(consts.StatusOK, report)
}

func (c *FraudReportController) GetFraudReportByTransactionHandler(ctx context.Context, reqCtx *app.RequestContext) {
	idStr := reqCtx.Param("transaction_id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		reqCtx.JSON(consts.StatusBadRequest, utils.H{
			"code":    consts.StatusBadRequest,
			"message": "Bad Request",
			"details": "Invalid transaction ID",
		})
		return
	}
	report, err := c.fraudReportService.GetByTransactionID(uint(id))
	if err != nil {
		reqCtx.JSON(consts.StatusNotFound, utils.H{
			"code":    consts.StatusNotFound,
			"message": "Not Found",
			"details": "Fraud report not found for this transaction",
		})
		return
	}
	reqCtx.JSON(consts.StatusOK, report)
}

func (c *FraudReportController) ListFraudReportsHandler(ctx context.Context, reqCtx *app.RequestContext) {
	pageStr := reqCtx.Query("page")
	sizeStr := reqCtx.Query("size")
	page, err := strconv.Atoi(pageStr)
	if err != nil || page <= 0 {
		page = 1
	}
	size, err := strconv.Atoi(sizeStr)
	if err != nil || size <= 0 {
		size = 10
	}
	reports, err := c.fraudReportService.List(page, size)
	if err != nil {
		reqCtx.JSON(consts.StatusInternalServerError, utils.H{
			"code":    consts.StatusInternalServerError,
			"message": "Internal Server Error",
			"details": err.Error(),
		})
		return
	}
	reqCtx.JSON(consts.StatusOK, reports)
}

func (c *FraudReportController) CreateFraudReportHandler(ctx context.Context, reqCtx *app.RequestContext) {
	var req service.FraudReportRequest
	if err := reqCtx.BindJSON(&req); err != nil {
		reqCtx.JSON(consts.StatusBadRequest, utils.H{
			"code":    consts.StatusBadRequest,
			"message": "Bad Request",
			"details": "Invalid request body",
		})
		return
	}
	id, err := c.fraudReportService.Create(&req)
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

func (c *FraudReportController) UpdateFraudReportHandler(ctx context.Context, reqCtx *app.RequestContext) {
	idStr := reqCtx.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		reqCtx.JSON(consts.StatusBadRequest, utils.H{
			"code":    consts.StatusBadRequest,
			"message": "Bad Request",
			"details": "Invalid fraud report ID",
		})
		return
	}
	var req service.FraudReportRequest
	if err := reqCtx.BindJSON(&req); err != nil {
		reqCtx.JSON(consts.StatusBadRequest, utils.H{
			"code":    consts.StatusBadRequest,
			"message": "Bad Request",
			"details": "Invalid request body",
		})
		return
	}
	report, err := c.fraudReportService.Update(uint(id), &req)
	if err != nil {
		reqCtx.JSON(consts.StatusInternalServerError, utils.H{
			"code":    consts.StatusInternalServerError,
			"message": "Internal Server Error",
			"details": err.Error(),
		})
		return
	}
	reqCtx.JSON(consts.StatusOK, report)
}

func (c *FraudReportController) DeleteFraudReportHandler(ctx context.Context, reqCtx *app.RequestContext) {
	idStr := reqCtx.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		reqCtx.JSON(consts.StatusBadRequest, utils.H{
			"code":    consts.StatusBadRequest,
			"message": "Bad Request",
			"details": "Invalid fraud report ID",
		})
		return
	}
	err = c.fraudReportService.Delete(uint(id))
	if err != nil {
		reqCtx.JSON(consts.StatusInternalServerError, utils.H{
			"code":    consts.StatusInternalServerError,
			"message": "Internal Server Error",
			"details": err.Error(),
		})
		return
	}
	reqCtx.Status(consts.StatusNoContent)
}

func (c *FraudReportController) GenerateReportHandler(ctx context.Context, reqCtx *app.RequestContext) {
	idStr := reqCtx.Param("transaction_id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		reqCtx.JSON(consts.StatusBadRequest, utils.H{
			"code":    consts.StatusBadRequest,
			"message": "Bad Request",
			"details": "Invalid transaction ID",
		})
		return
	}
	report, err := c.fraudReportService.GenerateReport(uint(id))
	if err != nil {
		reqCtx.JSON(consts.StatusInternalServerError, utils.H{
			"code":    consts.StatusInternalServerError,
			"message": "Internal Server Error",
			"details": err.Error(),
		})
		return
	}
	reqCtx.JSON(consts.StatusOK, report)
}
