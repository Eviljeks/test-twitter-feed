package app

import (
	"github.com/Eviljeks/test-twitter-feed/internal/api/route/sse"
	"github.com/Eviljeks/test-twitter-feed/internal/subscriber"
	"github.com/gin-gonic/gin"
)

func NewHandler(cfg *Config, agent *subscriber.Agent) (*gin.Engine, error) {
	newHandler := sse.NewNewHandler(cfg.NewMessageEventName, agent)

	r := gin.Default()

	api := r.Group("/sse")

	newHandler.Handle(api)

	return r, nil
}
