package app

import (
	"github.com/Eviljeks/test-twitter-feed/internal/amqp/publisher"
	"github.com/Eviljeks/test-twitter-feed/internal/api/route/message"
	"github.com/Eviljeks/test-twitter-feed/internal/store"
	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v4"
)

func NewHandler(cfg *Config, conn *pgx.Conn, publisher *publisher.SingleQueueAMQPPublisher) (*gin.Engine, error) {
	s := store.NewStore(conn)

	listHandler, err := message.NewListHandler(s)
	if err != nil {
		return nil, err
	}

	addHandler := message.NewAddHandler(s, publisher)

	r := gin.Default()

	api := r.Group("/api")

	addHandler.Handle(api)

	listHandler.Handle(r)

	return r, nil
}
