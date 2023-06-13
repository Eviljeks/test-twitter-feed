package app

import (
	"github.com/Eviljeks/test-twitter-feed/internal/api/route.go/message"
	"github.com/Eviljeks/test-twitter-feed/internal/messages"
	"github.com/Eviljeks/test-twitter-feed/internal/store"
	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v4"
)

func NewHandler(cfg *Config, conn *pgx.Conn) (*gin.Engine, error) {
	s := store.NewStore(conn)

	addedCh := make(chan messages.Message, cfg.AddedChanCap)
	listHandler, err := message.NewListHandler(s, addedCh)
	if err != nil {
		return nil, err
	}

	addHandler := message.NewAddHandler(s, addedCh)
	newHandler := message.NewNewHandler(cfg.NewMessageEventName, addedCh)

	r := gin.Default()

	api := r.Group("/api")

	addHandler.Handle(api)
	newHandler.Handle(api)

	listHandler.Handle(r)

	return r, nil
}
