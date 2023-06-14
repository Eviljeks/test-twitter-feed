package bot

import (
	"os"

	"github.com/Eviljeks/test-twitter-feed/cmd/bot/app"
	"github.com/spf13/cobra"
)

func NewBotCommand() *cobra.Command {
	cfg := app.NewConfig(os.Getenv("API_BASE_PATH"))

	cmd := &cobra.Command{
		Use: "bot",
		Run: func(cmd *cobra.Command, args []string) {

			cfg.Run()
		},
	}

	cmd.PersistentFlags().UintVar(&cfg.DelaySec, "delay", uint(0), "delay sec")
	cmd.PersistentFlags().UintVar(&cfg.RequestsPerMin, "reqs", uint(60), "requests per min")

	return cmd
}
