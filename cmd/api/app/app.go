package app

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/jackc/pgx/v4"
	"github.com/sirupsen/logrus"
	"github.com/streadway/amqp"
	"golang.org/x/sync/errgroup"

	iamqp "github.com/Eviljeks/test-twitter-feed/internal/amqp"
	"github.com/Eviljeks/test-twitter-feed/internal/http/health"
	"github.com/Eviljeks/test-twitter-feed/internal/store"
	"github.com/Eviljeks/test-twitter-feed/pkg/amqputil"
	"github.com/Eviljeks/test-twitter-feed/pkg/pgutil"
)

type Config struct {
	Port              string
	HealthPort        string
	MessagesQueueName string
}

func NewConfig(messagesQueueName string) *Config {
	return &Config{
		MessagesQueueName: messagesQueueName,
		Port:              ":3000",
		HealthPort:        ":5000",
	}
}

func (c *Config) Run() {
	var (
		databaseURL = os.Getenv("DATABASE_URL")
		amqpURL     = os.Getenv("AMQP_URL")
	)

	ctx := context.Background()

	eg := &errgroup.Group{}

	var conn *pgx.Conn
	var amqpConn *amqp.Connection

	eg.Go(func() error {
		var err error
		// setup db
		conn, err = pgutil.ConnectWithWait(ctx, databaseURL, time.Second, uint8(10))
		if err != nil {
			return fmt.Errorf("db connection failed, err: %s", err.Error())
		}
		logrus.Infoln("db connected")

		return nil
	})

	eg.Go(func() error {
		// setup amqp
		var err error
		amqpConn, err = amqputil.Connect(ctx, amqpURL, time.Second, uint8(10))
		if err != nil {
			return fmt.Errorf("amqp connect failed, err: %s", err.Error())
		}

		return nil
	})

	err := eg.Wait()
	if err != nil {
		if conn != nil {
			conn.Close(ctx)
		}

		if amqpConn != nil {
			amqpConn.Close()
		}

		panic(err)
	}

	defer conn.Close(ctx)
	defer amqpConn.Close()

	ch, err := amqpConn.Channel()
	if err != nil {
		panic(fmt.Sprintf("amqp channel failed, err: %s", err.Error()))
	}
	defer ch.Close()

	err = iamqp.Setup(ch, c.MessagesQueueName)
	if err != nil {
		panic(fmt.Sprintf("amqp setup failed, err: %s", err.Error()))
	}

	publisher := iamqp.NewSingleQueueAMQPPublisher(ch, c.MessagesQueueName)
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
