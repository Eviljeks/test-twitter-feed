package app

import (
	"os"
	"os/signal"
	"syscall"

	"github.com/sirupsen/logrus"
)

type Config struct {
	NewMessageEventName string
	MessagesQueueName   string
	AddedChanCap        int
	Port                string
}

func DefaultConfig(messagesQueueName string) *Config {
	return &Config{
		NewMessageEventName: "messages",
		MessagesQueueName:   messagesQueueName,
		AddedChanCap:        10,
		Port:                ":3000",
	}
}

func (c *Config) Run() {
	r, err := NewHandler(c)
	if err != nil {
		panic(err)
	}

	go func() {
		sErr := r.Run(c.Port)
		if sErr != nil {
			logrus.Fatalf("failed to run server: %v", sErr)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt, syscall.SIGTERM)

	<-quit

	logrus.Print("Server received shutdown signal")
}
