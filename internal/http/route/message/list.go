package message

import (
	"io"
	"net/http"

	"github.com/Eviljeks/test-twitter-feed/internal/messages"
	"github.com/Eviljeks/test-twitter-feed/internal/store"
	"github.com/gin-gonic/gin"
)

type Renderer interface {
	Render(tmplName string, data interface{}, w io.Writer) error
}

type ListHandler struct {
	store    *store.Store
	renderer Renderer
	ssePath  string
}

func NewListHandler(store *store.Store, renderer Renderer, ssePath string) (*ListHandler, error) {
	return &ListHandler{
		store:    store,
		renderer: renderer,
		ssePath:  ssePath,
	}, nil
}

func (lh *ListHandler) Handle(ctx *gin.Context) {
	ctx.Header("Content-Type", "text/html; charset=utf-8")

	msgs, err := lh.store.ListMessages(ctx)
	if err != nil {
		err := lh.renderer.Render("error.tmpl.html", nil, ctx.Writer)
		if err != nil {
			ctx.AbortWithStatus(http.StatusInternalServerError)
		}

		return
	}

	data := struct {
		Msgs    []*messages.Message
		SSEPath string
	}{
		Msgs:    msgs,
		SSEPath: lh.ssePath,
	}

	err = lh.renderer.Render("feed.tmpl.html", data, ctx.Writer)
	if err != nil {
		err = lh.renderer.Render("error.tmpl.html", nil, ctx.Writer)
		if err != nil {
			ctx.AbortWithStatus(http.StatusInternalServerError)
		}
	}
}
