package app

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"path/filepath"
	"syscall"
	"time"

	"github.com/sirupsen/logrus"

	"github.com/Eviljeks/test-twitter-feed/internal/amqp"
	"github.com/Eviljeks/test-twitter-feed/internal/store"
	"github.com/Eviljeks/test-twitter-feed/internal/templating"
	"github.com/Eviljeks/test-twitter-feed/pkg/amqputil"
	"github.com/Eviljeks/test-twitter-feed/pkg/pgutil"
)

const templatesRelativePath = "./../../../templates/"

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
	var (
		databaseURL = os.Getenv("DATABASE_URL")
		amqpURL     = os.Getenv("AMQP_URL")
		ssePath     = os.Getenv("SSE_PATH")
	)

	ctx := context.Background()

	// setup db
	conn, err := pgutil.ConnectWithWait(ctx, databaseURL, time.Second, uint8(10))
	if err != nil {
		panic(fmt.Sprintf("db connection failed, err: %s", err.Error()))
	}
	logrus.Infoln("db connected")
	defer conn.Close(ctx)

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

	err = amqp.Setup(ch, c.MessagesQueueName)
	if err != nil {
		panic(fmt.Sprintf("amqp setup failed, err: %s", err.Error()))
	}

	publisher := amqp.NewSingleQueueAMQPPublisher(ch, c.MessagesQueueName)
	if err != nil {
		panic(fmt.Sprintf("publisher failed, err: %s", err.Error()))
	}

	templatesBasePath, err := filepath.Abs(templatesRelativePath)
	if err != nil {
		panic(fmt.Sprintf("templates: %s", err.Error()))
	}

	s := store.NewStore(conn)
	renderer := templating.NewRenderer(templatesBasePath)

	// end setup

	r, err := NewHandler(publisher, s, renderer, ssePath)
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
