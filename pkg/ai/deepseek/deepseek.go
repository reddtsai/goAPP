package deepseek

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"

	"github.com/reddtsai/goAPP/pkg/ai"
)

type DeepSeekOption func(*DeepSeekOptions)

type DeepSeekOptions struct {
	BaseURL   string
	Header    http.Header
	Model     string
	MaxTokens int
}

func DefaultOptions() *DeepSeekOptions {
	return &DeepSeekOptions{
		BaseURL: "https://api.deepseek.com",
		Header: http.Header{
			"Content-Type": []string{"application/json"},
			"Accept":       []string{"application/json"},
		},
		Model: "deepseek-chat",
	}
}

func WithBaseURL(baseURL string) DeepSeekOption {
	return func(o *DeepSeekOptions) {
		o.BaseURL = baseURL
	}
}

func WithAuthorizationToken(token string) DeepSeekOption {
	return func(o *DeepSeekOptions) {
		o.Header.Add("Authorization", fmt.Sprintf("Bearer %s", token))
	}
}

type DeepSeekClient struct {
	options *DeepSeekOptions
	client  *http.Client
}

func NewDeepSeekClient(opts ...DeepSeekOption) *DeepSeekClient {
	options := DefaultOptions()
	for _, opt := range opts {
		opt(options)
	}

	client := &http.Client{
		Timeout: 10 * time.Second,
	}

	return &DeepSeekClient{
		options: options,
		client:  client,
	}
}

type CompletionRequest struct {
	Messages         []Message      `json:"messages"`
	Model            string         `json:"model"`
	FrequencyPenalty float64        `json:"frequency_penalty"`
	MaxTokens        int            `json:"max_tokens"`
	PresencePenalty  float64        `json:"presence_penalty"`
	ResponseFormat   ResponseFormat `json:"response_format"`
	// Stop             []string         `json:"stop"`
	Stream bool `json:"stream"`
	// StreamOptions    *StreamOptions  `json:"stream_options"`
	Temperature float32 `json:"temperature"`
	TopP        float32 `json:"top_p"`
	// Tools            interface{}     `json:"tools"` // null 或其他型別（可根據實際內容定義）
	// ToolChoice       string          `json:"tool_choice"`
	// Logprobs         bool            `json:"logprobs"`
	// TopLogprobs      *int            `json:"top_logprobs"` // null 或 int
}

type Message struct {
	Content string `json:"content"`
	Role    string `json:"role"`
}

type ResponseFormat struct {
	Type string `json:"type"`
}

type StreamOptions struct {
}

func (c *DeepSeekClient) Completion(params ai.CompletionParams) (*ai.CompletionResult, error) {
	// {
	// 	"messages": [
	// 	  {
	// 		"content": "You are a helpful assistant",
	// 		"role": "system"
	// 	  },
	// 	  {
	// 		"content": "Hi",
	// 		"role": "user"
	// 	  }
	// 	],
	// 	"stop": null,
	// 	"stream_options": null,
	// 	"tools": null,
	// 	"tool_choice": "none",
	// 	"logprobs": false,
	// 	"top_logprobs": null
	// }
	completion := &CompletionRequest{
		Model:            c.options.Model,
		FrequencyPenalty: 0,
		MaxTokens:        c.options.MaxTokens,
		PresencePenalty:  0,
		ResponseFormat: ResponseFormat{
			Type: "text",
		},
		Stream:      false,
		Temperature: 1,
		TopP:        1,
	}
	for _, msg := range params.Messages {
		completion.Messages = append(completion.Messages, Message{
			Content: msg.Content,
			Role:    msg.Role,
		})
	}
	jsonBytes, err := json.Marshal(completion)
	if err != nil {
		return nil, err
	}
	payload := strings.NewReader(string(jsonBytes))
	url := fmt.Sprintf("%s/chat/completions", c.options.BaseURL)
	req, err := http.NewRequest("POST", url, payload)
	if err != nil {
		return nil, err
	}
	req.Header = c.options.Header
	resp, err := c.client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	fmt.Println(string(body))

	result := &ai.CompletionResult{}
	return result, nil
}
