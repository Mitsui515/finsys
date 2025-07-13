// service/qwen.go
package service

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"io"
	"net/http"

	"github.com/Mitsui515/finsys/model"
)

type QwenLLM struct {
	apiKey       string
	client       *http.Client
	systemPrompt string
}

func NewQwenLLM(apiKey string) *QwenLLM {
	systemPrompt := `
# Role and Goal
You are FinBot, an advanced AI financial assistant embedded within the 'finsys' financial analytics platform. Your core duties are:
1.  **Data Analysis**: Analyze financial transaction data with precision.
2.  **Fraud Detection**: Identify and report potential fraudulent activities based on user requests.
3.  **Question Answering**: Address user inquiries regarding finance and transactions.
Your responses must always be professional, accurate, concise, and strictly adhere to security guidelines.

# Output Format Constraint
You MUST always respond with a JSON object containing a "thought" and a "response" key.
- **thought**: Briefly describe your thought process and the steps you plan to take to answer the user's query.
- **response**: Provide the final answer to the user. If the answer is text, provide a direct string. If it requires using a tool, provide a tool-call JSON object.

# Available Tools
When you need to retrieve information from the system, you MUST NOT fabricate data. You must request to call one of the following tools:
1.  **getTransactionByID**: Retrieves detailed information for a single transaction.
	- Format: ` + "`" + `{"tool_name": "getTransactionByID", "arguments": {"id": <transaction_id>}}` + "`" + `
2.  **listTransactions**: Lists multiple transactions, with optional filters.
	- Format: ` + "`" + `{"tool_name": "listTransactions", "arguments": {"page": <page_number>, "size": <items_per_page>, "type": "<CASH_IN|CASH_OUT|etc>"}}` + "`" + `
3.  **generateFraudReport**: Generates a detailed fraud analysis report for a specified transaction.
	- Format: ` + "`" + `{"tool_name": "generateFraudReport", "arguments": {"transaction_id": <transaction_id>}}` + "`" + `

# Few-shot Examples

---
**Example 1: Knowledge-based Question**
[USER]: "What are the common types of financial fraud?"
[AI]:
{
  "thought": "The user is asking a general financial knowledge question. I do not need to use a tool and can answer directly from my knowledge base.",
  "response": "Common types of financial fraud include identity theft, phishing scams, investment schemes, and credit card fraud. Within the finsys platform, we primarily focus on transactional fraud like unauthorized transfers."
}
---
**Example 2: Request for Specific Transaction Data (Tool Call)**
[USER]: "Can you look up the details for transaction ID 10258?"
[AI]:
{
  "thought": "The user needs detailed information for a specific transaction. I should use the getTransactionByID tool to retrieve this data.",
  "response": {
    "tool_name": "getTransactionByID",
    "arguments": {
      "id": 10258
    }
  }
}
---
**Example 3: Request for Fraud Analysis (Tool Call)**
[USER]: "Please analyze transaction 77889 for any fraud risk."
[AI]:
{
  "thought": "The user is requesting a fraud analysis for a specific transaction. The most appropriate tool for this is generateFraudReport.",
  "response": {
    "tool_name": "generateFraudReport",
    "arguments": {
      "transaction_id": 77889
    }
  }
}
---
`
	return &QwenLLM{
		apiKey:       apiKey,
		client:       &http.Client{},
		systemPrompt: systemPrompt,
	}
}

type qwenAPIRequest struct {
	Model       model.Model        `json:"model"`
	Messages    []model.LLMMessage `json:"messages"`
	Temperature float64            `json:"temperature,omitempty"`
	TopP        float64            `json:"top_p,omitempty"`
}

// --- UPDATED: Response struct for OpenAI-compatible endpoint ---
type qwenAPIResponse struct {
	Choices []struct {
		Message struct {
			Content string `json:"content"`
		} `json:"message"`
	} `json:"choices"`
}

func (q *QwenLLM) Chat(ctx context.Context, req *model.LLMChatRequest) (*model.LLMChatResponse, error) {
	if q.apiKey == "" || q.apiKey == "YOUR_DASHSCOPE_API_KEY" {
		return nil, errors.New("Qwen API key is not configured in config/config.go")
	}

	apiURL := "https://dashscope.aliyuncs.com/compatible-mode/v1/chat/completions"

	messagesWithSystemPrompt := append([]model.LLMMessage{{Role: "system", Content: q.systemPrompt}}, req.Messages...)

	apiRequest := qwenAPIRequest{
		Model:       req.Model,
		Messages:    messagesWithSystemPrompt,
		Temperature: 0.7,
		TopP:        0.8,
	}

	reqBody, err := json.Marshal(apiRequest)
	if err != nil {
		return nil, err
	}

	httpReq, err := http.NewRequestWithContext(ctx, "POST", apiURL, bytes.NewBuffer(reqBody))
	if err != nil {
		return nil, err
	}

	httpReq.Header.Set("Authorization", "Bearer "+q.apiKey)
	httpReq.Header.Set("Content-Type", "application/json")

	resp, err := q.client.Do(httpReq)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode != http.StatusOK {
		return nil, errors.New("Qwen API request failed with status " + resp.Status + ": " + string(respBody))
	}

	var apiResponse qwenAPIResponse
	if err := json.Unmarshal(respBody, &apiResponse); err != nil {
		return nil, err
	}

	if len(apiResponse.Choices) == 0 || apiResponse.Choices[0].Message.Content == "" {
		return nil, errors.New("received an empty or invalid response from Qwen API")
	}

	return &model.LLMChatResponse{
		Content: apiResponse.Choices[0].Message.Content,
	}, nil
}
