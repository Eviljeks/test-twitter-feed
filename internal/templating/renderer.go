package templating

import (
	"io"
	"path/filepath"
	"text/template"
)

type Renderer struct {
	basePath string
}

func NewRenderer(basePath string) *Renderer {
	return &Renderer{
		basePath: basePath,
	}
}

func (r *Renderer) Render(tmplName string, data interface{}, w io.Writer) error {
	path := filepath.Join(r.basePath, tmplName)
	t, err := template.ParseFiles(path)

	if err != nil {
		return err
	}

	return t.Execute(w, data)
}
