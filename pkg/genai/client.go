package genai

import (
	"context"
)

const (
	ROLE_USER      = "user"
	ROLE_ASSISTANT = "assistant"
)

type CompletionParams struct {
	// General parameters
	MaxOutputTokens int
	Messages        []Message `json:"messages"`

	// OpenAI specific parameters

	// DeepSeek specific parameters

	// Anthropic specific parameters

	// Bedrock specific parameters
}

type CompletionResult struct {
	Message Message    `json:"message"`
	Usage   TokenUsage `json:"usage"`
}

type Message struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

type TokenUsage struct {
	InputTokens  int64 `json:"input_tokens"`
	OutputTokens int64 `json:"output_tokens"`
	TotalTokens  int64 `json:"total_tokens"`
}

type Client interface {
	Completion(ctx context.Context, params CompletionParams) (*CompletionResult, error)
}
