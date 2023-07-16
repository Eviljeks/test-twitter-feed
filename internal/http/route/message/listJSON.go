package message

import (
	"net/http"

	api "github.com/Eviljeks/test-twitter-feed/internal/http"
	"github.com/Eviljeks/test-twitter-feed/internal/store"
	"github.com/gin-gonic/gin"
)

type ListJSONHandler struct {
	store *store.Store
}

func NewListJSONHandler(store *store.Store) *ListJSONHandler {
	return &ListJSONHandler{
		store: store,
	}
}

func (ljh *ListJSONHandler) Handle(ctx *gin.Context) {
	msgs, err := ljh.store.ListMessages(ctx)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, api.ErrorBadRequest(err))
		return
	}

	ctx.JSON(http.StatusOK, api.OK(msgs))
}
