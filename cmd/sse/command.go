package sse

import (
	"github.com/Eviljeks/test-twitter-feed/cmd/sse/app"
	"github.com/spf13/cobra"
)

func NewSSECommand(messagesQueueName string) *cobra.Command {
	cmd := &cobra.Command{
		Use: "sse",
		Run: func(cmd *cobra.Command, args []string) {
			cfg := app.DefaultConfig(messagesQueueName)

			cfg.Run()
		},
	}

	return cmd
}
