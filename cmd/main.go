package main

import (
	"github.com/Eviljeks/test-twitter-feed/cmd/api"
	"github.com/Eviljeks/test-twitter-feed/cmd/bot"
	"github.com/Eviljeks/test-twitter-feed/cmd/sse"
	"github.com/spf13/cobra"
)

const queueName = "messages"

func main() {
	cmd := &cobra.Command{
		Use: "twitter-feed",
		Run: func(cmd *cobra.Command, args []string) {},
	}

	cmd.AddCommand(sse.NewSSECommand(queueName))
	cmd.AddCommand(api.NewServerCommand(queueName))
	cmd.AddCommand(bot.NewBotCommand())

	err := cmd.Execute()
	if err != nil {
		panic(err)
	}
}
