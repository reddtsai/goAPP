package main

import (
	"context"
	"fmt"
	"log"
	"net/http"

	"github.com/spf13/cobra"

	"github.com/reddtsai/goAPP/cmd"
	"github.com/reddtsai/goAPP/internal/microservices/chatbot"
	awsSDK "github.com/reddtsai/goAPP/pkg/aws/v2"
	transportHTTP "github.com/reddtsai/goAPP/pkg/transport/http"
)

var (
	CONFIG_PATH = "./cmd/chatbot" // --build-arg CONFIG_PATH=

	ROOT_COMMAND = &cobra.Command{
		Use:   "chatbot",
		Short: "Chatbot application",
	}
)

func main() {
	ctx := context.Background()
	srv, err := newChatbotHttpService(ctx)
	if err != nil {
		log.Fatalf("Error creating chatbot service: %v", err)
	}
	ROOT_COMMAND.AddCommand(cmd.NewHTTPCommand(srv))
	if err := ROOT_COMMAND.Execute(); err != nil {
		log.Fatalf("Error executing command: %v", err)
	}
}

func newChatbotHttpService(ctx context.Context) (http.Handler, error) {
	cfg, err := LoadConfig(CONFIG_PATH)
	if err != nil {
		return nil, err
	}

	router := transportHTTP.NewRouter(ctx)
	switch cfg.GenAIVendor {
	case "bedrock":
		aesCfg, err := awsSDK.LoadConfigWithAKSK(ctx, cfg.AwsBedrockRegion, cfg.AwsBedrockAccessKeyID, cfg.AwsBedrockSecretAccessKey)
		if err != nil {
			return nil, err
		}
		handler := chatbot.NewHttpServicWithBedrock(ctx, aesCfg, cfg.AwsBedrockModel)
		handler.RegisterRoutes(router)
	case "openai":
		// TODO: Implement OpenAI handler
	case "deepseek":
		// TODO: Implement DeepSeek handler
	default:
		return nil, fmt.Errorf("unsupported GenAI vendor: %s", cfg.GenAIVendor)
	}
	return router, nil
}
