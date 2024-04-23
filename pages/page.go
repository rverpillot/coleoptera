package pages

import (
	"embed"
	"io/fs"
	"log"
	"os"
)

//go:embed templates
var EmbedTemplatesFs embed.FS

//go:embed statics
var EmbedResourcesFs embed.FS

var (
	TemplatesFs fs.FS = EmbedTemplatesFs
	ResourcesFs fs.FS = EmbedResourcesFs
)

func SetDebugMode() {
	log.Print("Loading pages from filesystem")
	TemplatesFs = os.DirFS("pages")
	ResourcesFs = os.DirFS("pages")
}
