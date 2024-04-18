package pages

import (
	"embed"
	"path"

	"github.com/rverpillot/ihui/templating"
)

//go:embed templates
var TemplatesFs embed.FS

//go:embed statics
var ResourcesFs embed.FS

func newAceTemplate(name string, model interface{}) *templating.PageAce {
	return templating.NewPageAce(TemplatesFs, path.Join("templates", name), model)
}
