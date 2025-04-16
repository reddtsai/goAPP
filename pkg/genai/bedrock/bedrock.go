package bedrock

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/bedrockruntime"
	bt "github.com/aws/aws-sdk-go-v2/service/bedrockruntime/types"

	"github.com/reddtsai/goAPP/pkg/genai"
)

type BedrockOption func(*BedrockOptions)

type BedrockOptions struct {
	AwsCfg       aws.Config
	Model        string
	Temperature  float32
	SystemPrompt string
}

func DefaultOptions() *BedrockOptions {
	return &BedrockOptions{
		Temperature: 0.7,
	}
}

func WithAwsConfig(cfg aws.Config) BedrockOption {
	return func(o *BedrockOptions) {
		o.AwsCfg = cfg
	}
}

func WithModel(model string) BedrockOption {
	return func(o *BedrockOptions) {
		o.Model = model
	}
}

func WithTemperature(temp float32) BedrockOption {
	return func(o *BedrockOptions) {
		o.Temperature = temp
	}
}

func WithSystemPrompt(prompt string) BedrockOption {
	return func(o *BedrockOptions) {
		o.SystemPrompt = prompt
	}
}

type BedrockClient struct {
	options *BedrockOptions
	client  *bedrockruntime.Client
}

func NewBedrockClient(opts ...BedrockOption) *BedrockClient {
	options := DefaultOptions()
	for _, opt := range opts {
		opt(options)
	}

	c := bedrockruntime.NewFromConfig(options.AwsCfg)

	return &BedrockClient{
		options: options,
		client:  c,
	}
}

func (c *BedrockClient) Completion(ctx context.Context, params genai.CompletionParams) (*genai.CompletionResult, error) {
	return c.converse(ctx, params)
}

func (c *BedrockClient) converse(ctx context.Context, params genai.CompletionParams) (*genai.CompletionResult, error) {
	converseInput := &bedrockruntime.ConverseInput{
		ModelId: aws.String(c.options.Model),
		InferenceConfig: &bt.InferenceConfiguration{
			Temperature: aws.Float32(c.options.Temperature),
		},
	}
	converseInput.System = append(converseInput.System, &bt.SystemContentBlockMemberText{
		Value: c.options.SystemPrompt,
	})

	for _, message := range params.Messages {
		switch message.Role {
		case genai.ROLE_USER:
			userPrompt := bt.Message{
				Role: bt.ConversationRoleUser,
				Content: []bt.ContentBlock{
					&bt.ContentBlockMemberText{
						Value: message.Content,
					},
				},
			}
			converseInput.Messages = append(converseInput.Messages, userPrompt)
		case genai.ROLE_ASSISTANT:
			assistantPrompt := bt.Message{
				Role: bt.ConversationRoleAssistant,
				Content: []bt.ContentBlock{
					&bt.ContentBlockMemberText{
						Value: message.Content,
					},
				},
			}
			converseInput.Messages = append(converseInput.Messages, assistantPrompt)
		}
	}

	output, err := c.client.Converse(ctx, converseInput)
	if err != nil {
		return nil, err
	}

	// jsonResp, _ := json.MarshalIndent(output, "", "  ")
	// fmt.Println(string(jsonResp))
	result := &genai.CompletionResult{
		Usage: genai.TokenUsage{
			InputTokens:  int64(*output.Usage.InputTokens),
			OutputTokens: int64(*output.Usage.OutputTokens),
			TotalTokens:  int64(*output.Usage.TotalTokens),
		},
	}

	assistantMsg := output.Output.(*bt.ConverseOutputMemberMessage).Value
	if len(assistantMsg.Content) > 0 {
		result.Message = genai.Message{
			Role:    genai.ROLE_ASSISTANT,
			Content: assistantMsg.Content[0].(*bt.ContentBlockMemberText).Value,
		}
	}

	return result, nil
}
