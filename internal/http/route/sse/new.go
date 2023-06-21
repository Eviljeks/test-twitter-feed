package sse

import (
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

func (nh *NewHandler) Handle(ctx *gin.Context) {
	ch := nh.agent.Subscribe()
	defer logrus.Println("unsubscribed")
	defer nh.agent.Unsubscribe(ch)

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

				return false
			}
		}
	})
}
