package app

import (
	"context"

	"github.com/Eviljeks/test-twitter-feed/internal/api/route/sse"
	"github.com/Eviljeks/test-twitter-feed/internal/subscriber"
	"github.com/gin-gonic/gin"
)

func NewHandler(ctx context.Context, cfg *Config, agent *subscriber.Agent) (*gin.Engine, error) {
	newHandler := sse.NewNewHandler(cfg.NewMessageEventName, agent)

	r := gin.Default()

	api := r.Group("/sse")

	newHandler.Handle(ctx, api)

	return r, nil
}
