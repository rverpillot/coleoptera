package pages

import (
	"rverpi/ihui.v2"

	rice "github.com/GeertJohan/go.rice"
)

var (
	templateBox  = rice.MustFindBox("templates")
	ResourcesBox = rice.MustFindBox("statics")
)

func newAceTemplate(name string, model interface{}) *ihui.AceTemplateDrawer {
	content, err := templateBox.Bytes(name)
	if err != nil {
		panic(err)
	}
	return ihui.NewAceTemplateDrawer(content, model)
}

type Page struct {
	tmpl *ihui.AceTemplateDrawer
}

func NewPage(pageTemplate string, model interface{}) *Page {
	return &Page{
		templ: ihui.newAceTemplate(pageTemplate, model),
	}
}

func (p *Page) Draw(page ihui.PageDrawer) {
	page.Draw(p.tmpl)
}
