package message

import (
	"net/http"

	"github.com/Eviljeks/test-twitter-feed/internal/api"
	"github.com/Eviljeks/test-twitter-feed/internal/messages"
	"github.com/Eviljeks/test-twitter-feed/internal/store"
	"github.com/Eviljeks/test-twitter-feed/pkg/clock"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type Publisher interface {
	Publish(data interface{}) error
}

type AddHandler struct {
	store     *store.Store
	publisher Publisher
}

func NewAddHandler(store *store.Store, publisher Publisher) *AddHandler {
	return &AddHandler{
		store:     store,
		publisher: publisher,
	}
}

type AddParams struct {
	Content string `form:"content" json:"content"`
}

func (ah *AddHandler) Handle(r gin.IRouter) {
	r.POST("/messages", func(ctx *gin.Context) {
		var params AddParams

		if err := ctx.ShouldBind(&params); err != nil {
			ctx.JSON(http.StatusBadRequest, api.ErrorBadRequest(err))
			return
		}

		m := messages.Message{
			UUID:    uuid.NewString(),
			Content: params.Content,
			TS:      clock.GetCurrentTS(),
		}

		_, err := ah.store.SaveMessage(ctx, m)
		if err != nil {
			ctx.JSON(http.StatusBadRequest, api.ErrorInternalServer(err))

			return
		}

		err = ah.publisher.Publish(m)
		if err != nil {
			ctx.JSON(http.StatusBadRequest, api.ErrorBadRequest(err))

			return
		}

		ctx.JSON(http.StatusOK, api.OK(m))
	})
}
