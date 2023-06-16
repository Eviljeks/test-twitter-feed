package app

import (
	publisher "github.com/Eviljeks/test-twitter-feed/internal/amqp"
	"github.com/Eviljeks/test-twitter-feed/internal/api/route/message"
	"github.com/Eviljeks/test-twitter-feed/internal/store"
	"github.com/Eviljeks/test-twitter-feed/internal/templating"
	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v4"
)

func NewHandler(
	cfg *Config,
	conn *pgx.Conn,
	publisher *publisher.SingleQueueAMQPPublisher,
	ssePath string,
	templatesBasePath string,
) (*gin.Engine, error) {
	s := store.NewStore(conn)
	renderer := templating.NewRenderer(templatesBasePath)

	listHandler, err := message.NewListHandler(s, renderer, ssePath)
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
