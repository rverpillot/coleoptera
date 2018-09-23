package pages

import (
	"github.com/rverpillot/ihui"

	rice "github.com/GeertJohan/go.rice"
)

var (
	templateBox  = rice.MustFindBox("templates")
	ResourcesBox = rice.MustFindBox("statics")
)

func newAceTemplate(name string, model interface{}) *ihui.PageAce {
	content, err := templateBox.Bytes(name)
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
