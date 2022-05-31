package pages

import (
	"embed"
	"path"

	"github.com/rverpillot/ihui"
)

//go:embed templates
var TemplatesFs embed.FS

//go:embed statics
var ResourcesFs embed.FS

func newAceTemplate(name string, model interface{}) *ihui.PageAce {
	content, err := TemplatesFs.ReadFile(path.Join("templates", name))
	if err != nil {
		panic(err)
	}
	return ihui.NewPageAce(name, content, model)
}

type Page struct {
	tmpl *ihui.PageAce
}

func NewPage(pageTemplate string, model interface{}) *Page {
	return &Page{
		tmpl: newAceTemplate(pageTemplate, model),
	}
}

func (p *Page) Render(page ihui.Page) {
	p.tmpl.Render(page)
}
