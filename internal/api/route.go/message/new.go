package message

import (
	"io"

	"github.com/Eviljeks/test-twitter-feed/internal/messages"
	"github.com/gin-gonic/gin"
)

type NewHandler struct {
	eventName string
	addedCh   <-chan messages.Message
}

func NewNewHandler(eventName string, addedCh <-chan messages.Message) *NewHandler {
	return &NewHandler{
		eventName: eventName,
		addedCh:   addedCh,
	}
}

func (nh *NewHandler) Handle(r gin.IRouter) {
	r.GET("/messages/new", func(ctx *gin.Context) {
		ctx.Stream(func(w io.Writer) bool {
			if msg, ok := <-nh.addedCh; ok {
				ctx.SSEvent(nh.eventName, msg)
				return true
			}
			return false
		})
	})
}
