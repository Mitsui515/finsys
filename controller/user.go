package controller

import (
	"context"

	"github.com/Mitsui515/finsys/service"
	"github.com/cloudwego/hertz/pkg/app"
	"github.com/cloudwego/hertz/pkg/common/utils"
	"github.com/cloudwego/hertz/pkg/protocol/consts"
)

type UserController struct {
	userService *service.UserService
}

func NewUserController() *UserController {
	return &UserController{
		userService: service.NewUserService(),
	}
}

func (c *UserController) Register(ctx context.Context, reqCtx *app.RequestContext) {
	var req service.RegisterRequest
	if err := reqCtx.BindAndValidate(&req); err != nil {
		reqCtx.JSON(consts.StatusBadRequest, utils.H{
			"code":    consts.StatusBadRequest,
			"message": "Invalid Input",
			"details": err.Error(),
		})
		return
	}
	userID, err := c.userService.Register(&req)
	if err != nil {
		reqCtx.JSON(consts.StatusBadRequest, utils.H{
			"code":    consts.StatusBadRequest,
			"message": "Invalid Input",
			"details": err.Error(),
		})
	}
	reqCtx.JSON(consts.StatusCreated, utils.H{
		"user_id": userID,
	})
}

func (c *UserController) Login(ctx context.Context, reqCtx *app.RequestContext) {
	var req service.LoginRequest
	if err := reqCtx.BindAndValidate(&req); err != nil {
		reqCtx.JSON(consts.StatusBadRequest, utils.H{
			"code":    consts.StatusBadRequest,
			"message": "Invalid Input",
			"details": err.Error(),
		})
		return
	}
	token, err := c.userService.Login(&req)
	if err != nil {
		reqCtx.JSON(consts.StatusUnauthorized, utils.H{
			"code":    consts.StatusUnauthorized,
			"message": "Unauthorized",
			"details": "Invalid username or password",
		})
		return
	}
	reqCtx.JSON(consts.StatusOK, utils.H{
		"token": token,
	})
}

func (c *UserController) GetUserInfo(ctx context.Context, reqCtx *app.RequestContext) {
	userID, exists := reqCtx.Get("user_id")
	if !exists {
		reqCtx.JSON(consts.StatusForbidden, utils.H{
			"code":    consts.StatusForbidden,
			"message": "Forbidden",
			"details": "Invalid or expired token",
		})
		return
	}
	user, err := c.userService.GetByID(userID.(uint))
	if err != nil {
		reqCtx.JSON(consts.StatusNotFound, utils.H{
			"code":    consts.StatusNotFound,
			"message": "Not Found",
			"details": err.Error(),
		})
		return
	}
	reqCtx.JSON(consts.StatusOK, utils.H{
		"user_id":    userID,
		"username":   user.Username,
		"email":      user.Email,
		"created_at": user.CreatedAt,
	})
}
