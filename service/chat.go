// service/chat.go
package service

import (
	"context"
	"errors"
	"fmt"
	"strings"

	"github.com/Mitsui515/finsys/config"
	"github.com/Mitsui515/finsys/model"
)

type ChatService struct {
	llms map[string]model.LLM
}

type ChatRequest struct {
	Prompt string      `json:"prompt" binding:"required"`
	Model  model.Model `json:"model"`
}

type ChatResponse struct {
	Reply string `json:"reply"`
}

func NewChatService() *ChatService {
	appConfig := config.DefaultConfig()
	llms := make(map[string]model.LLM)
	llms["qwen"] = NewQwenLLM(appConfig.LLM.APIKey)
	return &ChatService{
		llms: llms,
	}
}

func (s *ChatService) HandleChat(ctx context.Context, req *ChatRequest) (*ChatResponse, error) {
	if req.Prompt == "" {
		return nil, errors.New("prompt cannot be empty")
	}
	modelName := req.Model
	if modelName == "" {
		modelName = model.QwenTurbo
	}
	var llmClient model.LLM
	if strings.HasPrefix(string(modelName), "qwen") {
		llmClient = s.llms["qwen"]
	}
	if llmClient == nil {
		return nil, fmt.Errorf("unsupported model: %s", modelName)
	}
	llmReq := &model.LLMChatRequest{
		Model: modelName,
		Messages: []model.LLMMessage{
			{Role: "user", Content: req.Prompt},
		},
	}
	llmResp, err := llmClient.Chat(ctx, llmReq)
	if err != nil {
		return nil, fmt.Errorf("failed to get reply from LLM: %w", err)
	}
	return &ChatResponse{Reply: llmResp.Content}, nil
}
