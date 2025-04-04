package main

import (
	"log"

	"github.com/spf13/cobra"

	"github.com/reddtsai/goAPP/cmd"
)

var RootCmd = &cobra.Command{
	Use:   "chatbot",
	Short: "Chatbot application",
}

func init() {
	RootCmd.AddCommand(cmd.HttpCmd)
}

func main() {
	if err := RootCmd.Execute(); err != nil {
		log.Fatalf("Error executing command: %v", err)
	}
}
