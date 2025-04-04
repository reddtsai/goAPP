package ai

type CompletionParams struct {
	// General parameters
	Messages []Message `json:"messages"`

	// OpenAI specific parameters

	// DeepSeek specific parameters

	// Anthropic specific parameters

	// Bedrock specific parameters
}

type CompletionResult struct {
}

type Message struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

type Client interface {
	Completion(CompletionParams) (CompletionResult, error)
}
