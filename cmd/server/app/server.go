package app

import (
	"context"
	"os"
	"os/signal"
	"syscall"

	"github.com/sirupsen/logrus"

	"github.com/Eviljeks/test-twitter-feed/pkg/pgutil"
)

type Config struct {
	NewMessageEventName string
	AddedChanCap        int
}

func DefaultConfig() *Config {
	return &Config{
		NewMessageEventName: "messages",
		AddedChanCap:        10,
	}
}

func (c *Config) Run() {
	ctx := context.Background()
	conn, err := pgutil.Connect(ctx, os.Getenv("DATABASE_URL"))
	if err != nil {
		logrus.Fatalf("db connection faild, err: %s", err.Error())
	}
	logrus.Infoln("db connected")
	defer conn.Close(ctx)

	r, err := NewHandler(c, conn)
	if err != nil {
		panic(err)
	}

	go func() {
		sErr := r.Run(":3000")
		if sErr != nil {
			logrus.Fatalf("failed to run server: %v", sErr)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt, syscall.SIGTERM)

	<-quit

	logrus.Print("Server received shutdown signal")
}
