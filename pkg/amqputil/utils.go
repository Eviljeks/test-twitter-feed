package amqputil

import (
	"context"
	"strings"
	"time"

	"github.com/Eviljeks/test-twitter-feed/pkg"
	"github.com/streadway/amqp"
)

func Connect(ctx context.Context, connString string, tick time.Duration, numTicks uint8) (*amqp.Connection, error) {
	var conn *amqp.Connection

	err := pkg.NewWaiter(tick, numTicks).Wait(ctx, func(ctx context.Context) error {
		var err error
		conn, err = amqp.Dial(connString)
		if err != nil {
			if strings.Contains(err.Error(), "dial tcp") {
				return pkg.ErrNotReadyYet
			}

			return err
		}

		if err != nil {
			return err
		}
		return nil
	})
	if err != nil {
		return nil, err
	}

	return conn, nil
}
