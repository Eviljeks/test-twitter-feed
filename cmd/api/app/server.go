package app

import (
	"time"

	"github.com/Eviljeks/test-twitter-feed/internal/amqp"
	"github.com/Eviljeks/test-twitter-feed/internal/http/route/message"
	"github.com/Eviljeks/test-twitter-feed/internal/store"
	"github.com/Eviljeks/test-twitter-feed/internal/templating"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func NewServer(
	publisher *amqp.SingleQueueAMQPPublisher,
	store *store.Store,
	renderer *templating.Renderer,
	ssePath string,
) (*gin.Engine, error) {
	listHTMLHandler, err := message.NewListHTMLHandler(store, renderer, ssePath)
	if err != nil {
		return nil, err
	}

	listJSONHandler := message.NewListJSONHandler(store)

	addHandler := message.NewAddHandler(store, publisher)

	r := gin.Default()

	r.Use(cors.New(cors.Config{
		AllowOrigins:  []string{"*"},
		AllowMethods:  []string{"GET"},
		AllowHeaders:  []string{"Origin"},
		ExposeHeaders: []string{"Content-Length"},
		MaxAge:        12 * time.Hour,
	}))

	r.GET("/messages", listHTMLHandler.Handle)

	g := r.Group("/api")

	g.POST("/message", addHandler.Handle)
	g.GET("/messages", listJSONHandler.Handle)

	return r, nil
}
