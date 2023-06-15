package sse

import (
	"context"
	"io"

	"github.com/Eviljeks/test-twitter-feed/internal/subscriber"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

type NewHandler struct {
	eventName string
	agent     *subscriber.Agent
}

func NewNewHandler(eventName string, agent *subscriber.Agent) *NewHandler {
	return &NewHandler{
		eventName: eventName,
		agent:     agent,
	}
}

func (nh *NewHandler) Handle(ctx context.Context, r gin.IRouter) {
	r.GET("/messages/new", func(ctx *gin.Context) {
		ch := nh.agent.Subscribe()

		ctx.Header("Access-Control-Allow-Origin", "*")

		ctx.Stream(func(w io.Writer) bool {
			for {
				select {
				case msg, ok := <-ch:
					if ok {
						ctx.SSEvent(nh.eventName, msg)
						return true
					}

					return false
				case <-ctx.Done():
					logrus.Print("SSE ctx.Done")
					return false
				}
			}
		})
	})
}
