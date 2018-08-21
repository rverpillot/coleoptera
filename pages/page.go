package pages

import (
	"io"
	"log"
	"rverpi/ihui"

	rice "github.com/GeertJohan/go.rice"
	"github.com/yosssi/ace"
)

var (
	templateBox  = rice.MustFindBox("templates")
	ResourcesBox = rice.MustFindBox("statics")
)

func renderTemplate(pageTemplate string, w io.Writer, model interface{}) {
	opts := &ace.Options{DynamicReload: true}
	opts.Asset = func(name string) ([]byte, error) {
		data, err := templateBox.Bytes(name)
		return data, err
	}
	tmpl, err := ace.Load(pageTemplate, "", opts)
	if err != nil {
		log.Println(err)
		panic(err)
	}

	err = tmpl.Execute(w, model)
	if err != nil {
		log.Println(err)
	}
}

type Page struct {
	*ihui.Page
	pageTemplate string
}

func NewPage(pageTemplate string, modal bool) *Page {
	page := &Page{
		Page:         ihui.NewPage("Coleoptera"),
		pageTemplate: pageTemplate,
	}

	return page
}

func (p *Page) renderPage(w io.Writer, model interface{}) {
	renderTemplate(p.pageTemplate, w, model)
}
