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
	"github.com/Eviljeks/test-twitter-feed/internal/http/health"
	"github.com/Eviljeks/test-twitter-feed/internal/store"
	"github.com/Eviljeks/test-twitter-feed/pkg/amqputil"
	"github.com/Eviljeks/test-twitter-feed/pkg/pgutil"
)

type Config struct {
	Port                string
	HealthPort          string
	MessagesQueueName   string
	PostgresTimeoutSecs uint8
	AMQPTimeoutSecs     uint8
}

func NewConfig(messagesQueueName string) *Config {
	return &Config{
		MessagesQueueName:   messagesQueueName,
		Port:                ":3000",
		HealthPort:          ":5000",
		PostgresTimeoutSecs: 30,
		AMQPTimeoutSecs:     30,
	}
}

func (c *Config) Run() {
	var (
		databaseURL = os.Getenv("DATABASE_URL")
		amqpURL     = os.Getenv("AMQP_URL")
	)

	ctx := context.Background()

	// setup db
	conn, err := pgutil.ConnectWithWait(ctx, databaseURL, time.Second, c.PostgresTimeoutSecs)
	if err != nil {
		panic(fmt.Sprintf("db connection failed, err: %s", err.Error()))
	}
	logrus.Infoln("db connected")
	defer conn.Close(ctx)

	// setup amqp
	amqpConn, err := amqputil.Connect(ctx, amqpURL, time.Second, c.AMQPTimeoutSecs)
	if err != nil {
		panic(fmt.Sprintf("amqp connect failed, err: %s", err.Error()))
	}
	defer amqpConn.Close()

	ch, err := amqpConn.Channel()
	if err != nil {
		panic(fmt.Sprintf("amqp channel failed, err: %s", err.Error()))
	}
	defer ch.Close()

	err = amqp.Setup(ch, c.MessagesQueueName)
	if err != nil {
		panic(fmt.Sprintf("amqp setup failed, err: %s", err.Error()))
	}

	publisher := amqp.NewSingleQueueAMQPPublisher(ch, c.MessagesQueueName)
	if err != nil {
		panic(fmt.Sprintf("publisher failed, err: %s", err.Error()))
	}

	s := store.NewStore(conn)

	// end setup

	server, err := NewServer(publisher, s)
	if err != nil {
		panic(err)
	}

	go func() {
		sErr := server.Run(c.Port)
		if sErr != nil {
			logrus.Fatalf("failed to run server: %v", sErr)
		}
	}()

	h := health.NewServer()

	go func() {
		sErr := h.Run(c.HealthPort)
		if sErr != nil {
			logrus.Fatalf("failed to metrics server: %v", sErr)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt, syscall.SIGTERM)

	<-quit

	logrus.Print("Api server received shutdown signal")
}
