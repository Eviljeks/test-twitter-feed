package app

import (
	"github.com/Eviljeks/test-twitter-feed/internal/amqp"
	"github.com/Eviljeks/test-twitter-feed/internal/api/route/message"
	"github.com/Eviljeks/test-twitter-feed/internal/store"
	"github.com/Eviljeks/test-twitter-feed/internal/templating"
	"github.com/gin-gonic/gin"
)

func NewHandler(
	publisher *amqp.SingleQueueAMQPPublisher,
	store *store.Store,
	renderer *templating.Renderer,
	ssePath string,
) (*gin.Engine, error) {
	listHandler, err := message.NewListHandler(store, renderer, ssePath)
	if err != nil {
		return nil, err
	}

	addHandler := message.NewAddHandler(store, publisher)

	r := gin.Default()

	api := r.Group("/api")

	addHandler.Handle(api)

	listHandler.Handle(r)

	return r, nil
}
