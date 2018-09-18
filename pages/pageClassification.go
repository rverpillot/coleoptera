package pages

import (
	"log"

	"bitbucket.org/rverpi90/ihui"

	"github.com/jinzhu/gorm"
	"rverpi90/coleoptera.v3/model"
)

type PageClassification struct {
	tmpl           *ihui.PageAce
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

func (page *PageClassification) Render(p ihui.Page) {
	db := p.Get("db").(*gorm.DB)

	page.tmpl.Render(p)

	p.On("click", "close", func(s *ihui.Session, event ihui.Event) bool {
		return s.QuitPage()
	})

	p.On("submit", "form", func(s *ihui.Session, event ihui.Event) bool {
		data := event.Data.(map[string]interface{})
		page.classification.Nom = data["classification"].(string)
		if err := db.Create(page.classification).Error; err != nil {
			log.Println(err)
			page.Error = err.Error()
		} else {
			s.QuitPage()
		}
		return true
	})
}
