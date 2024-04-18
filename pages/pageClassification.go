package pages

import (
	"log"

	"github.com/rverpillot/ihui"
	"github.com/rverpillot/ihui/templating"

	"github.com/jinzhu/gorm"
	"github.com/rverpillot/coleoptera/model"
)

type PageClassification struct {
	tmpl           *templating.PageAce
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

func (page *PageClassification) Render(p *ihui.Page) error {
	db := p.Get("db").(*gorm.DB)

	if err := page.tmpl.Render(p); err != nil {
		return err
	}

	p.On("click", "close", func(s *ihui.Session, event ihui.Event) {
		p.Close()
	})

	p.On("submit", "form", func(s *ihui.Session, event ihui.Event) {
		data := event.Data.(map[string]interface{})
		page.classification.Nom = data["classification"].(string)
		if err := db.Create(page.classification).Error; err != nil {
			log.Println(err)
			page.Error = err.Error()
		} else {
			p.Close()
		}
	})

	return nil
}
