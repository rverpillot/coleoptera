package pages

import (
	"github.com/rverpillot/ihui"

	"github.com/rverpillot/coleoptera/model"
	"gorm.io/gorm"
)

type PageClassification struct {
	classification *model.Classification
	Error          string
}

func newPageClassification(classification *model.Classification) *PageClassification {
	return &PageClassification{
		classification: classification,
	}
}

func (page *PageClassification) Render(p *ihui.Page) error {
	db := p.Get("db").(*gorm.DB)

	if err := p.WriteGoTemplate(TemplatesFs, "templates/classification.html", page); err != nil {
		return err
	}

	p.On("click", "#close", func(s *ihui.Session, event ihui.Event) error {
		return p.Close()
	})

	p.On("submit", "form", func(s *ihui.Session, event ihui.Event) error {
		data := event.Data.(map[string]interface{})
		page.classification.Nom = data["classification"].(string)
		if page.classification.Nom == "" {
			page.Error = "Nom obligatoire"
			return nil
		}
		if err := db.Create(page.classification).Error; err != nil {
			page.Error = err.Error()
			return nil
		}
		return p.Close()
	})

	return nil
}
