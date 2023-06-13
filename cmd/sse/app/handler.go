package app

import (
	"github.com/Eviljeks/test-twitter-feed/internal/api/route/sse"
	"github.com/Eviljeks/test-twitter-feed/internal/messages"
	"github.com/gin-gonic/gin"
)

func NewHandler(cfg *Config) (*gin.Engine, error) {
	addedCh := make(chan messages.Message, cfg.AddedChanCap)
	newHandler := sse.NewNewHandler(cfg.NewMessageEventName, addedCh)

	r := gin.Default()

	api := r.Group("/sse")

	newHandler.Handle(api)

	return r, nil
}
