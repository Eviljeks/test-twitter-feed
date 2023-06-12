package store

import (
	"context"

	"github.com/pkg/errors"

	"github.com/Eviljeks/test-twitter-feed/internal/messages"
	"github.com/Eviljeks/test-twitter-feed/pkg/pgutil"
)

const (
	MessagesTable         = "messages"
	MessagesColumnUUID    = "uuid"
	MessagesColumnContent = "content"
	MessagesColumnTS      = "ts"
)

func (s *Store) SaveMessage(ctx context.Context, m messages.Message) (bool, error) {
	sb := pgutil.SB()
	sql, args, err := sb.Insert(MessagesTable).
		Columns(
			MessagesColumnUUID,
			MessagesColumnContent,
			MessagesColumnTS,
		).
		Values(
			m.UUID,
			m.Content,
			m.TS,
		).
		ToSql()

	if err != nil {
		return false, errors.Wrapf(err, "message[%s]", m.UUID)
	}

	tag, err := s.Conn.Exec(ctx, sql, args...)
	if err != nil {
		return false, errors.Wrapf(err, "message[%s]", m.UUID)
	}

	return tag.RowsAffected() > 0, nil
}

func (s *Store) ListMessages(ctx context.Context) ([]*messages.Message, error) {
	sql, args, err := pgutil.SB().
		Select(
			MessagesColumnUUID,
			MessagesColumnContent,
			MessagesColumnTS,
		).
		From(MessagesTable).
		OrderBy(MessagesColumnTS + " DESC").
		ToSql()

	if err != nil {
		return nil, err
	}

	rows, err := s.Conn.Query(ctx, sql, args...)
	if err != nil {
		return nil, err
	}

	results := make([]*messages.Message, 0)

	for rows.Next() {
		var m messages.Message

		if err = rows.Scan(
			&m.UUID,
			&m.Content,
			&m.TS,
		); err != nil {
			return nil, err
		}

		results = append(results, &m)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return results, nil
}
