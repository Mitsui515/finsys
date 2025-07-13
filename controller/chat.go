package controller

import (
	"context"

	"github.com/Mitsui515/finsys/service"
	"github.com/cloudwego/hertz/pkg/app"
	"github.com/cloudwego/hertz/pkg/protocol/consts"
)

type ChatController struct {
	chatService *service.ChatService
}

func NewChatController() *ChatController {
	return &ChatController{
		chatService: service.NewChatService(),
	}
}

func (c *ChatController) ChatHandler(ctx context.Context, reqCtx *app.RequestContext) {
	var req service.ChatRequest
	if err := reqCtx.BindAndValidate(&req); err != nil {
		reqCtx.JSON(consts.StatusBadRequest, map[string]string{"error": "Invalid request"})
		return
	}
	resp, err := c.chatService.HandleChat(ctx, &req)
	if err != nil {
		reqCtx.JSON(consts.StatusInternalServerError, map[string]string{"error": err.Error()})
		return
	}
	reqCtx.JSON(consts.StatusOK, resp)
}
