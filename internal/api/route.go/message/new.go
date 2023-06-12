package message

import (
	"io"

	"github.com/Eviljeks/test-twitter-feed/internal/messages"
	"github.com/Eviljeks/test-twitter-feed/internal/store"
	"github.com/gin-gonic/gin"
)

type NewHandler struct {
	addedCh <-chan messages.Message
}

func NewNewHandler(store *store.Store, addedCh <-chan messages.Message) *NewHandler {
	return &NewHandler{
		addedCh: addedCh,
	}
}

func (nh *NewHandler) Handle(r gin.IRouter) {
	r.GET("/messages", func(ctx *gin.Context) {
		ctx.Stream(func(w io.Writer) bool {
			if msg, ok := <-nh.addedCh; ok {
				ctx.JSON(200, msg)
				return true
			}
			return false
		})
	})
}
