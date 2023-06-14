package api

import (
	"github.com/Eviljeks/test-twitter-feed/cmd/api/app"
	"github.com/spf13/cobra"
)

func NewServerCommand(queueName string) *cobra.Command {
	cmd := &cobra.Command{
		Use: "api",
		Run: func(cmd *cobra.Command, args []string) {
			cfg := app.NewConfig(queueName)

			cfg.Run()
		},
	}

	return cmd
}
