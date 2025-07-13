package model

import "context"

type Model string

const (
	QwenTurbo Model = "qwen-turbo"
	GPT4o     Model = "gpt-4o"
)

type LLM interface {
	Chat(ctx context.Context, req *LLMChatRequest) (*LLMChatResponse, error)
}

type LLMMessage struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

type LLMChatRequest struct {
	Model    Model        `json:"model"`
	Messages []LLMMessage `json:"messages"`
}

type LLMChatResponse struct {
	Content string `json:"content"`
}
