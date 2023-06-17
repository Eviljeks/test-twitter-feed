package bot

import (
	"github.com/Eviljeks/test-twitter-feed/cmd/bot/app"
	"github.com/spf13/cobra"
)

func NewBotCommand() *cobra.Command {
	cfg := app.NewConfig()

	cmd := &cobra.Command{
		Use: "bot",
		Run: func(cmd *cobra.Command, args []string) {

			cfg.Run()
		},
	}

	cmd.PersistentFlags().UintVar(&cfg.DelaySec, "delay", uint(0), "delay sec")
	cmd.PersistentFlags().UintVar(&cfg.RequestsPerMin, "reqs-per-min", uint(60), "requests per min")

	return cmd
}
