package app

import (
	"context"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/Eviljeks/test-twitter-feed/internal/http/health"
	"github.com/Eviljeks/test-twitter-feed/pkg/client"
	"github.com/jaswdr/faker"
	"github.com/sirupsen/logrus"
)

type Config struct {
	DelaySec       uint
	RequestsPerMin uint
	HealthPort     string
}

func NewConfig() *Config {
	return &Config{
		HealthPort: ":5000",
	}
}

func (c *Config) Run() {
	var (
		apiBasePath = os.Getenv("API_BASE_PATH")
	)

	ctx := context.Background()

	httpClient := http.DefaultClient

	apiClient := client.NewAPIClient(httpClient, apiBasePath)

	faker := faker.New()

	bot := NewBot(c.DelaySec, c.RequestsPerMin, apiClient, &faker)

	go (func() {
		bot.Run(ctx)
	})()

	h := health.NewServer()

	go func() {
		sErr := h.Run(c.HealthPort)
		if sErr != nil {
			logrus.Fatalf("failed to health server: %v", sErr)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt, syscall.SIGTERM)

	<-quit

	logrus.Print("Bot received shutdown signal")
}
