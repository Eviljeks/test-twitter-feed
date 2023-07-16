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

type ListHTMLHandler struct {
	store    *store.Store
	renderer Renderer
	ssePath  string
}

func NewListHTMLHandler(store *store.Store, renderer Renderer, ssePath string) (*ListHTMLHandler, error) {
	return &ListHTMLHandler{
		store:    store,
		renderer: renderer,
		ssePath:  ssePath,
	}, nil
}

func (lhh *ListHTMLHandler) Handle(ctx *gin.Context) {
	ctx.Header("Content-Type", "text/html; charset=utf-8")

	msgs, err := lhh.store.ListMessages(ctx)
	if err != nil {
		err := lhh.renderer.Render("error.tmpl.html", nil, ctx.Writer)
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
		SSEPath: lhh.ssePath,
	}

	err = lhh.renderer.Render("feed.tmpl.html", data, ctx.Writer)
	if err != nil {
		err = lhh.renderer.Render("error.tmpl.html", nil, ctx.Writer)
		if err != nil {
			ctx.AbortWithStatus(http.StatusInternalServerError)
		}
	}
}
