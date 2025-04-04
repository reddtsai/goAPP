package cmd

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/spf13/cobra"

	"github.com/reddtsai/goAPP/internal/microservices/user"
	transportHttp "github.com/reddtsai/goAPP/pkg/transport/http"
)

func init() {
	HttpCmd.Flags().String("service", "", "The service to run (e.g. 'user')")
	_ = HttpCmd.MarkFlagRequired("service")
}

var HttpCmd = &cobra.Command{
	Use:   "http",
	Short: "Run HTTP server for a selected service",
	RunE: func(cmd *cobra.Command, args []string) error {
		serviceName, _ := cmd.Flags().GetString("service")
		if serviceName == "" {
			return fmt.Errorf("missing required flag: --service")
		}
		return runHTTP(context.Background(), serviceName)
	},
}

func runHTTP(ctx context.Context, serviceName string) error {
	router, err := LoadHttpService(ctx, serviceName)
	if err != nil {
		return fmt.Errorf("failed to load HTTP service: %w", err)
	}

	srv := &http.Server{
		Addr:    ":8080",
		Handler: router,
	}
	go func() {
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("ListenAndServe: %v", err)
		}
	}()

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, syscall.SIGTERM, syscall.SIGINT)
	<-stop
	shutdownCtx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	if err := srv.Shutdown(shutdownCtx); err != nil {
		return fmt.Errorf("graceful shutdown failed: %w", err)
	}

	return nil
}

func LoadHttpService(ctx context.Context, name string) (http.Handler, error) {
	router := transportHttp.NewRouter(ctx)

	switch name {
	case "user":
		userService := user.NewHttpHandler(ctx)
		userService.RegisterRoutes(router)
	default:
		return nil, fmt.Errorf("service %s not found", name)
	}

	return router, nil
}
