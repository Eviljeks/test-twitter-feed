package pgutil

import (
	"context"
	"strings"
	"time"

	"github.com/Eviljeks/test-twitter-feed/pkg"
	"github.com/Masterminds/squirrel"
	"github.com/jackc/pgx/v4"
)

func SB() squirrel.StatementBuilderType {
	return squirrel.StatementBuilder.PlaceholderFormat(squirrel.Dollar)
}

func connect(ctx context.Context, connString string) (*pgx.Conn, error) {
	var conn *pgx.Conn
	conn, err := pgx.Connect(ctx, connString)
	if err != nil {
		return nil, err
	}

	return conn, nil
}

func ConnectWithWait(ctx context.Context, connString string, tick time.Duration, numTicks uint8) (*pgx.Conn, error) {
	var conn *pgx.Conn
	err := pkg.NewWaiter(tick, numTicks).Wait(ctx, func(ctx context.Context) error {
		var err error
		conn, err = connect(ctx, connString)
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
