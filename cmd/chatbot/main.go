package main

import (
	"context"
	"log"
	"net/http"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	"github.com/reddtsai/goAPP/cmd"
	"github.com/reddtsai/goAPP/internal/microservices/chatbot"
	awsSDK "github.com/reddtsai/goAPP/pkg/aws/v2"
	transportHTTP "github.com/reddtsai/goAPP/pkg/transport/http"
)

var RootCmd = &cobra.Command{
	Use:   "chatbot",
	Short: "Chatbot application",
}

func main() {
	ctx := context.Background()
	srv, err := newChatbotHttpService(ctx)
	if err != nil {
		log.Fatalf("Error creating chatbot service: %v", err)
	}
	RootCmd.AddCommand(cmd.NewHTTPCommand(srv))
	if err := RootCmd.Execute(); err != nil {
		log.Fatalf("Error executing command: %v", err)
	}
}

func newChatbotHttpService(ctx context.Context) (http.Handler, error) {
	v := viper.New()
	v.AutomaticEnv()

	// TODO
	router := transportHTTP.NewRouter(ctx)
	aesCfg, err := awsSDK.LoadConfigWithAKSK(ctx, v.GetString("REGION"), v.GetString("AK"), v.GetString("SK"))
	if err != nil {
		return nil, err
	}
	handler := chatbot.NewHttpServicWithBedrock(ctx, aesCfg, v.GetString("MODEL"))
	handler.RegisterRoutes(router)
	return router, nil
}
