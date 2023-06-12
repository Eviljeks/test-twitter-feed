package main

import "github.com/Eviljeks/test-twitter-feed/cmd/server/app"

func main() {
	cfg := app.DefaultConfig()

	cfg.Run()
}
