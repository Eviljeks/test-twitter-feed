package app

import (
	"time"

	"github.com/Eviljeks/test-twitter-feed/internal/amqp"
	"github.com/Eviljeks/test-twitter-feed/internal/http/route/message"
	"github.com/Eviljeks/test-twitter-feed/internal/store"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func NewServer(
	publisher *amqp.SingleQueueAMQPPublisher,
	store *store.Store,
	frontendURL string,
) (*gin.Engine, error) {
	listJSONHandler := message.NewListHandler(store)

	addHandler := message.NewAddHandler(store, publisher)

	r := gin.Default()

	r.Use(cors.New(cors.Config{
		AllowOrigins:  []string{frontendURL},
		AllowMethods:  []string{"GET"},
		AllowHeaders:  []string{"Origin"},
		ExposeHeaders: []string{"Content-Length"},
		MaxAge:        12 * time.Hour,
	}))

	g := r.Group("/api")

	g.POST("/message", addHandler.Handle)
	g.GET("/messages", listJSONHandler.Handle)

	return r, nil
}
