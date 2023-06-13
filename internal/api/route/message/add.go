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

type AddHandler struct {
	store *store.Store
}

func NewAddHandler(store *store.Store) *AddHandler {
	return &AddHandler{
		store: store,
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

		ctx.JSON(http.StatusOK, api.OK(m))
	})
}
