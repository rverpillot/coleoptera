package pages

import (
	"log"

	"github.com/jinzhu/gorm"
	"rverpi/coleoptera.v3/model"
	"rverpi/ihui.v2"
)

type PageClassification struct {
	tmpl           *ihui.AceTemplateDrawer
	classification *model.Classification
	Error          string
}

func newPageClassification(classification *model.Classification) *PageClassification {
	page := &PageClassification{
		classification: classification,
	}
	page.tmpl = newAceTemplate("classification.ace", page)
	return page
}

func (page *PageClassification) Draw(p ihui.Page) {
	db := p.Get("db").(*gorm.DB)

	p.Draw(page.tmpl)

	p.On("click", "close", func(s *ihui.Session, event ihui.Event) {
		s.QuitPage()
	})

	p.On("submit", "form", func(s *ihui.Session, event ihui.Event) {
		data := event.Data.(map[string]interface{})
		page.classification.Nom = data["classification"].(string)
		if err := db.Create(page.classification).Error; err != nil {
			log.Println(err)
			page.Error = err.Error()
		} else {
			s.QuitPage()
		}
	})
}
