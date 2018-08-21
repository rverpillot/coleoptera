package pages

import (
	"io"
	"log"
	"rverpi/coleoptera/model"
	"rverpi/ihui"

	"github.com/jinzhu/gorm"
)

type PageClassification struct {
	*Page
	classification *model.Classification
	CloseAction    ihui.ClickAction
	SubmitAction   ihui.SubmitAction
	Error          string
}

func newPageClassification(classification *model.Classification) ihui.PageRender {
	page := &PageClassification{
		Page:           NewPage("classification", false),
		classification: classification,
	}
	return page
}

func (page *PageClassification) OnInit(ctx *ihui.Context) {
	db := ctx.Get("db").(*gorm.DB)

	page.CloseAction = func(_ string) {
		page.Close()
	}

	page.SubmitAction = func(data map[string]interface{}) {
		page.classification.Nom = data["classification"].(string)
		if err := db.Create(page.classification).Error; err != nil {
			log.Println(err)
			page.Error = err.Error()
		} else {
			page.Close()
		}
	}
}

func (page *PageClassification) Render(w io.Writer, ctx *ihui.Context) {
	//db := ctx.Get("db").(*gorm.DB)

	page.renderPage(w, page)

}
