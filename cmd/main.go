package main

import (
	"github.com/Eviljeks/test-twitter-feed/cmd/api"
	"github.com/Eviljeks/test-twitter-feed/cmd/sse"
	"github.com/spf13/cobra"
)

func main() {
	cmd := &cobra.Command{
		Use: "twitter-feed",
		Run: func(cmd *cobra.Command, args []string) {},
	}

	cmd.AddCommand(sse.NewSSECommand())
	cmd.AddCommand(api.NewServerCommand())

	err := cmd.Execute()
	if err != nil {
		panic(err)
	}
}
