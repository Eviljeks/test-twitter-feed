package sse

import (
	"github.com/Eviljeks/test-twitter-feed/cmd/sse/app"
	"github.com/spf13/cobra"
)

func NewSSECommand() *cobra.Command {
	cmd := &cobra.Command{
		Use: "sse",
		Run: func(cmd *cobra.Command, args []string) {
			cfg := app.DefaultConfig()

			cfg.Run()
		},
	}

	return cmd
}
