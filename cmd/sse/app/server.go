package app

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/sirupsen/logrus"

	"github.com/Eviljeks/test-twitter-feed/internal/amqp"
	"github.com/Eviljeks/test-twitter-feed/internal/subscriber"
	"github.com/Eviljeks/test-twitter-feed/pkg/amqputil"
)

type Config struct {
	NewMessageEventName string
	MessagesQueueName   string
	Port                string
	ShutdownDelay       uint8
}

func NewConfig(messagesQueueName string) *Config {
	return &Config{
		NewMessageEventName: "messages",
		MessagesQueueName:   messagesQueueName,
		Port:                ":3000",
		ShutdownDelay:       2,
	}
}

func (c *Config) Run() {
	var (
		amqpURL = os.Getenv("AMQP_URL")
	)

	ctx, cancel := context.WithCancel(context.Background())

	// setup amqp
	amqpConn, err := amqputil.Connect(ctx, amqpURL, time.Second, uint8(10))
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
		for {
			select {
			case d := <-msgs:
				agent.Publish(string(d.Body))
			case <-ctx.Done():
				ch.Cancel(c.MessagesQueueName, false)
				return

			}
		}
	}()

	r, err := NewHandler(ctx, c, agent)
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

	logrus.Print("SSE server received shutdown signal")

	cancel()

	time.Sleep(time.Second * time.Duration(c.ShutdownDelay))

	logrus.Print("SSE canceled")
}
