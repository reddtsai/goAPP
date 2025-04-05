package chatbot

import (
	"context"
	"fmt"

	"github.com/aws/aws-sdk-go-v2/aws"

	"github.com/reddtsai/goAPP/pkg/genai"
	"github.com/reddtsai/goAPP/pkg/genai/bedrock"
)

type Service struct {
	ctx             context.Context
	generativeAiApi genai.Client

	GenerativeAiModel            string
	GenerativeAiApiContextWindow int
	GenerativeAiMaxOutputTokens  int
	GenerativeAiTemperature      float32
}

func NewHttpServicWithBedrock(ctx context.Context, awsConfig aws.Config, genaiModel string) *HttpHandler {
	srv := newServiceWithBedrock(ctx, awsConfig, genaiModel)
	h := &HttpHandler{
		service: srv,
	}

	return h
}

func newServiceWithBedrock(ctx context.Context, awsConfig aws.Config, genaiModel string) *Service {
	generativeAi := bedrock.NewBedrockClient(
		bedrock.WithAwsConfig(awsConfig),
		bedrock.WithModel(genaiModel),
		bedrock.WithTemperature(0.7),
	)
	return &Service{
		ctx:             ctx,
		generativeAiApi: generativeAi,
	}
}

type SendMessageParams struct {
	Message string `json:"message"`
}

type SendMessageResult struct {
	Answer string `json:"answer"`
}

func (s *Service) SendMessage(params SendMessageParams) (*SendMessageResult, error) {
	completion, err := s.generativeAiApi.Completion(s.ctx, genai.CompletionParams{
		Messages: []genai.Message{
			{
				Role:    genai.ROLE_USER,
				Content: params.Message,
			},
		},
	})
	if err != nil {
		fmt.Println("Error:", err)
		return nil, err
	}

	return &SendMessageResult{
		Answer: completion.Message.Content,
	}, nil
}
