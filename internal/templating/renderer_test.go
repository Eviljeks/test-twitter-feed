package templating

import (
	"bytes"
	"path/filepath"
	"testing"
	"time"

	"github.com/Eviljeks/test-twitter-feed/internal/messages"
)

func TestRenderer_Render(t *testing.T) {
	now := time.Now().UnixMicro()
	templatesBasePath, _ := filepath.Abs("./../../templates/")

	type fields struct {
		basePath string
	}
	type args struct {
		tmplName string
		data     interface{}
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantW   string
		wantErr bool
	}{
		{
			name: "",
			fields: fields{
				basePath: templatesBasePath,
			},
			args: args{
				tmplName: "feed.tmpl.html",
				data: struct {
					Msgs    []*messages.Message
					SSEPath string
				}{
					Msgs: []*messages.Message{
						{
							UUID:    "adae1b56-702d-4833-9618-342825ee7920",
							Content: "Some fresh message",
							TS:      now,
						},
					},
					SSEPath: "http://sse-path.com",
				},
			},
			wantW: `<!DOCTYPE html>
<html>
	<head>
		<meta charset="UTF-8">
		<title>Twitter feed</title>
	</head>
	<body>
		<h3>Feed</h3>
		<div id="feed">
		<div>adae1b56-702d-4833-9618-342825ee7920: Some fresh message </div>
		</div>
	</body>
	<script>
	const evtSource = new EventSource("http://sse-path.com");

	evtSource.addEventListener("messages", (event) => {
		const newElement = document.createElement("div");
		const eventList = document.getElementById("feed");
		const data = JSON.parse(event.data);
		newElement.textContent = data.uuid + ": " + data.content;
		eventList.prepend(newElement);
	  });
	</script>
</html>`,
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := &Renderer{
				basePath: tt.fields.basePath,
			}
			w := &bytes.Buffer{}
			if err := r.Render(tt.args.tmplName, tt.args.data, w); (err != nil) != tt.wantErr {
				t.Errorf("Renderer.Render() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if gotW := w.String(); gotW != tt.wantW {
				t.Errorf("Renderer.Render() = %v, want %v", gotW, tt.wantW)
			}
		})
	}
}
