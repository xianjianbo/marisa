package chat

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"

	"github.com/go-resty/resty/v2"
)

const (
	apiKey = "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"
)

var (
	RoleUser      string = "user"
	RoleAssistant string = "assistant"
	RoleSystem    string = "system"
)

var (
	ModelGPT35TURBO string = "gpt-3.5-turbo"
)

type Message struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

type GPTRequest struct {
	Model       string      `json:"model"`
	Messages    []Message   `json:"messages"`
	MaxTokens   int         `json:"max_tokens,omitempty"`
	Temperature int         `json:"temperature,omitempty"`
	TopP        int         `json:"top_p,omitempty"`
	N           int         `json:"n,omitempty"`
	Stream      bool        `json:"stream,omitempty"`
	Logprobs    interface{} `json:"logprobs,omitempty"`
	Stop        string      `json:"stop,omitempty"`
}

type GPTResponse struct {
	ID      string `json:"id"`
	Object  string `json:"object"`
	Created int    `json:"created"`
	Model   string `json:"model"`
	Usage   struct {
		PromptTokens     int `json:"prompt_tokens"`
		CompletionTokens int `json:"completion_tokens"`
		TotalTokens      int `json:"total_tokens"`
	} `json:"usage"`
	Choices []struct {
		Message      Message `json:"message"`
		FinishReason string  `json:"finish_reason"`
		Index        int     `json:"index"`
	} `json:"choices"`
}

func GPT(ctx context.Context, msgs []Message) (replyContent string, err error) {
	client := resty.New()
	clientResp, err := client.R().
		SetHeader("Content-Type", "application/json").
		SetHeader("Authorization", fmt.Sprintf("Bearer %s", apiKey)).
		SetBody(GPTRequest{
			Model:    ModelGPT35TURBO,
			Messages: msgs,
		}).
		Post("https://api.openai.com/v1/chat/completions")

	if err != nil {
		return
	}

	var resp GPTResponse
	if err = json.Unmarshal(clientResp.Body(), &resp); err != nil {
		return
	}
	if len(resp.Choices) == 0 {
		err = errors.New("length of choices is zero" + string(clientResp.Body()))
		return
	}

	replyContent = resp.Choices[0].Message.Content

	return
}
