package app

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/Eviljeks/test-twitter-feed/internal/amqp"
	"github.com/Eviljeks/test-twitter-feed/internal/subscriber"
	"github.com/Eviljeks/test-twitter-feed/pkg/amqputil"
	"github.com/sirupsen/logrus"
)

type Config struct {
	NewMessageEventName string
	MessagesQueueName   string
	Port                string
}

func DefaultConfig(messagesQueueName string) *Config {
	return &Config{
		NewMessageEventName: "messages",
		MessagesQueueName:   messagesQueueName,
		Port:                ":3000",
	}
}

func (c *Config) Run() {
	ctx := context.Background()

	// setup amqp
	amqpConn, err := amqputil.Connect(ctx, os.Getenv("AMQP_URL"), time.Second, uint8(5))
	if err != nil {
		panic(fmt.Sprintf("amqp connect failed, err: %s", err.Error()))
	}
	defer amqpConn.Close()

	ch, err := amqpConn.Channel()
	if err != nil {
		panic(fmt.Sprintf("amqp channel failed, err: %s", err.Error()))
	}
	defer ch.Close()

	amqp.Setup(ch, c.MessagesQueueName)

	msgs, err := ch.Consume(
		c.MessagesQueueName,
		"",
		true,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		panic(fmt.Sprintf("amqp consume failed, err: %s", err.Error()))
	}

	agent := subscriber.NewAgent()
	defer agent.Close()

	go func() {
		for d := range msgs {
			agent.Publish(string(d.Body)) // TODO move to consumer
		}
	}()

	r, err := NewHandler(c, agent)
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
