package pages

import (
	"embed"
)

//go:embed templates
var TemplatesFs embed.FS

//go:embed statics
var ResourcesFs embed.FS
