package pgutil

import (
	"context"

	"github.com/Masterminds/squirrel"
	"github.com/jackc/pgx/v4"
)

func SB() squirrel.StatementBuilderType {
	return squirrel.StatementBuilder.PlaceholderFormat(squirrel.Dollar)
}

func Connect(ctx context.Context, connString string) (*pgx.Conn, error) {
	var conn *pgx.Conn
	conn, err := pgx.Connect(ctx, connString)
	if err != nil {
		return nil, err
	}

	return conn, nil
}
