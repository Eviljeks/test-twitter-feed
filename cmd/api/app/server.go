package app

import (
	"github.com/Eviljeks/test-twitter-feed/internal/amqp"
	"github.com/Eviljeks/test-twitter-feed/internal/http/route/message"
	"github.com/Eviljeks/test-twitter-feed/internal/store"
	"github.com/Eviljeks/test-twitter-feed/internal/templating"
	"github.com/gin-gonic/gin"
)

func NewServer(
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

	r.POST("/api/messages", addHandler.Handle)
	r.GET("/messages", listHandler.Handle)

	return r, nil
}
