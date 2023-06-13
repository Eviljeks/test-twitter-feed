package message

import (
	"html/template"
	"net/http"

	"github.com/Eviljeks/test-twitter-feed/internal/api"
	"github.com/Eviljeks/test-twitter-feed/internal/messages"
	"github.com/Eviljeks/test-twitter-feed/internal/store"
	"github.com/gin-gonic/gin"
)

const tpl = `
<!DOCTYPE html>
<html>
	<head>
		<meta charset="UTF-8">
		<title>Twitter feed</title>
	</head>
	<body>
		<h3>Feed</h3>
		<div id="feed">
		{{range .Msgs}}<div>{{ .UUID }}: {{ .Content }} </div>{{end}}
		</div>
	</body>
	<script>
	const evtSource = new EventSource("http://localhost:3000/api/messages/new");

	evtSource.addEventListener("messages", (event) => {
		const newElement = document.createElement("div");
		const eventList = document.getElementById("feed");
		const data = JSON.parse(event.data);
		newElement.textContent = data.uuid + ": " + data.content;
		eventList.prepend(newElement);
	  });
	</script>
</html>`

type ListHandler struct {
	store    *store.Store
	addedCh  <-chan messages.Message
	template *template.Template
}

func NewListHandler(store *store.Store, addedCh <-chan messages.Message) (*ListHandler, error) {
	t, err := template.New("webpage").Parse(tpl)
	if err != nil {
		return nil, err
	}

	return &ListHandler{
		store:    store,
		addedCh:  addedCh,
		template: t,
	}, nil
}

func (lh *ListHandler) Handle(r gin.IRouter) {
	r.GET("/messages", func(ctx *gin.Context) {
		msgs, err := lh.store.ListMessages(ctx)
		if err != nil {
			ctx.JSON(http.StatusBadRequest, api.ErrorBadRequest(err))

			return
		}

		data := struct {
			Msgs []*messages.Message
		}{
			Msgs: msgs,
		}

		lh.template.Execute(ctx.Writer, data)
	})
}
