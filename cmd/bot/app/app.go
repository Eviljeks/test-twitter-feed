package app

import (
	"context"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/Eviljeks/test-twitter-feed/pkg/client"
	"github.com/jaswdr/faker"
	"github.com/sirupsen/logrus"
)

type Config struct {
	DelaySec       uint
	RequestsPerMin uint
	ApiBasePath    string
}

func NewConfig(apiBasePath string) *Config {
	return &Config{
		ApiBasePath: apiBasePath,
	}
}

func (c *Config) Run() {
	ctx := context.Background()

	httpClient := http.DefaultClient

	apiClient := client.NewApiClient(httpClient, c.ApiBasePath)

	faker := faker.New()

	bot := NewBot(c.DelaySec, c.RequestsPerMin, apiClient, &faker)

	go (func() {
		bot.Run(ctx)
	})()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt, syscall.SIGTERM)

	<-quit

	logrus.Print("Bot received shutdown signal")
}
