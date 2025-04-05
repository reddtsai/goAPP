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
)

var (
	Addr = ":8080"
)

func NewHTTPCommand(handler http.Handler) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "http",
		Short: "Run HTTP server for a selected service",
		RunE: func(cmd *cobra.Command, args []string) error {
			if handler == nil {
				return fmt.Errorf("handler is nil")
			}
			return runHttp(handler)
		},
	}

	cmd.Flags().StringVar(&Addr, "addr", ":8080", "Address to listen on (e.g., :8080)")
	// _ = cmd.MarkFlagRequired("addr竹手戈")

	return cmd
}

func runHttp(handler http.Handler) error {
	srv := &http.Server{
		Addr:    Addr,
		Handler: handler,
	}
	go func() {
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("ListenAndServe: %v", err)
		}
	}()

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, syscall.SIGTERM, syscall.SIGINT)
	<-stop
	shutdownCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := srv.Shutdown(shutdownCtx); err != nil {
		return fmt.Errorf("graceful shutdown failed: %w", err)
	}

	return nil
}
