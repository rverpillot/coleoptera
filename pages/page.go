package pages

import (
	"embed"
	"path"

	"github.com/rverpillot/ihui"
	"github.com/rverpillot/ihui/templating"
)

//go:embed templates
var TemplatesFs embed.FS

//go:embed statics
var ResourcesFs embed.FS

func newAceTemplate(name string) ihui.Template {
	return templating.NewPageFileAce(TemplatesFs, path.Join("templates", name))
}
