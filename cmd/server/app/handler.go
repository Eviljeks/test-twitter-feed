package app

import (
	"github.com/Eviljeks/test-twitter-feed/internal/api/route.go/message"
	"github.com/Eviljeks/test-twitter-feed/internal/messages"
	"github.com/Eviljeks/test-twitter-feed/internal/store"
	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v4"
)

func NewHandler(cfg *Config, conn *pgx.Conn) *gin.Engine {
	s := store.NewStore(conn)

	addedCh := make(chan messages.Message, cfg.AddedChanCap)

	addHandler := message.NewAddHandler(s, addedCh)
	listHandler := message.NewListHandler(s, addedCh)

	r := gin.Default()

	addHandler.Handle(r)
	listHandler.Handle(r)

	return r
}
