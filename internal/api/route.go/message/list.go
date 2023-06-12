package message

import (
	"io"
	"net/http"

	"github.com/Eviljeks/test-twitter-feed/internal/api"
	"github.com/Eviljeks/test-twitter-feed/internal/messages"
	"github.com/Eviljeks/test-twitter-feed/internal/store"
	"github.com/gin-gonic/gin"
)

type ListHandler struct {
	store   *store.Store
	addedCh <-chan messages.Message
}

func NewListHandler(store *store.Store, addedCh <-chan messages.Message) *ListHandler {
	return &ListHandler{
		store:   store,
		addedCh: addedCh,
	}
}

func (lh *ListHandler) Handle(r gin.IRouter) {
	r.GET("/messages", func(ctx *gin.Context) {
		msgs, err := lh.store.ListMessages(ctx)
		if err != nil {
			ctx.JSON(http.StatusBadRequest, api.ErrorBadRequest(err))

			return
		}

		for _, msg := range msgs {
			ctx.JSON(200, msg)
		}

		ctx.Stream(func(w io.Writer) bool {
			if msg, ok := <-lh.addedCh; ok {
				ctx.JSON(200, msg)
				return true
			}
			return false
		})
	})
}
