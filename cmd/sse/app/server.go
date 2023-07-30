package app

import (
	"context"
	"time"

	"github.com/Eviljeks/test-twitter-feed/internal/http/route/sse"
	"github.com/Eviljeks/test-twitter-feed/internal/subscriber"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func NewHandler(ctx context.Context, cfg *Config, agent *subscriber.Agent) (*gin.Engine, error) {
	newHandler := sse.NewNewHandler(cfg.NewMessageEventName, agent)

	r := gin.Default()

	r.Use(cors.New(cors.Config{
		AllowOrigins:  []string{"*"},
		AllowMethods:  []string{"GET"},
		AllowHeaders:  []string{"Origin"},
		ExposeHeaders: []string{"Content-Length"},
		MaxAge:        12 * time.Hour,
	}))

	r.GET("/sse/messages/new", newHandler.Handle)

	return r, nil
}
