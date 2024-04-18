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
	return ihui.NewPageAce(TemplatesFs, path.Join("templates", name), model)
}
