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
	"github.com/Eviljeks/test-twitter-feed/pkg/amqputil"
	"github.com/Eviljeks/test-twitter-feed/pkg/pgutil"
)

type Config struct {
	Port              string
	MessagesQueueName string
}

func NewConfig(messagesQueueName string) *Config {
	return &Config{
		MessagesQueueName: messagesQueueName,
		Port:              ":3000",
	}
}

func (c *Config) Run() {
	ctx := context.Background()

	// setup db
	conn, err := pgutil.Connect(ctx, os.Getenv("DATABASE_URL"))
	if err != nil {
		panic(fmt.Sprintf("db connection failed, err: %s", err.Error()))
	}
	logrus.Infoln("db connected")
	defer conn.Close(ctx)

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

	publisher := amqp.NewSingleQueueAMQPPublisher(ch, c.MessagesQueueName)
	if err != nil {
		panic(fmt.Sprintf("publisher failed, err: %s", err.Error()))
	}

	// end setup

	r, err := NewHandler(c, conn, publisher)
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

	logrus.Print("Api server received shutdown signal")
}
