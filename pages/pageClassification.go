package pages

import (
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
	return &PageClassification{
		tmpl:           newAceTemplate("classification.ace"),
		classification: classification,
	}
}

func (page *PageClassification) Render(p *ihui.Page) error {
	db := p.Get("db").(*gorm.DB)

	if err := p.WriteTemplate(page.tmpl, page); err != nil {
		return err
	}

	p.On("click", "close", func(s *ihui.Session, event ihui.Event) error {
		return p.Close()
	})

	p.On("submit", "form", func(s *ihui.Session, event ihui.Event) error {
		data := event.Data.(map[string]interface{})
		page.classification.Nom = data["classification"].(string)
		if err := db.Create(page.classification).Error; err != nil {
			page.Error = err.Error()
			return err
		} else {
			return p.Close()
		}
	})

	return nil
}
