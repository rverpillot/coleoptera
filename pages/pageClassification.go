package pages

import (
	"log"

	"github.com/rverpillot/ihui"

	"github.com/rverpillot/coleoptera/model"
	"gorm.io/gorm"
)

type PageClassification struct {
	Classification *model.Classification
	EspecesCount   int64
	Delete         bool
	Error          string
}

func newPageClassification(classification *model.Classification) *PageClassification {
	return &PageClassification{
		Classification: classification,
	}
}

func (page *PageClassification) Render(e *ihui.HTMLElement) error {
	db := e.Get("db").(*gorm.DB)

	e.OnClick("#close", func(s *ihui.Session, event ihui.Event) error {
		return e.Close()
	})

	e.OnClick("#delete", func(s *ihui.Session, ev ihui.Event) error {
		if err := db.Model(&model.Espece{}).Where("classification_id = ?", page.Classification.ID).Count(&page.EspecesCount).Error; err != nil {
			log.Println(err)
			page.Error = err.Error()
			return nil
		}
		page.Delete = true
		return nil
	})

	e.OnClick("#confirm-delete", func(s *ihui.Session, ev ihui.Event) error {
		var ids []uint
		for _, espece := range page.Classification.Especes {
			ids = append(ids, espece.ID)
		}
		if len(ids) > 0 {
			if err := db.Where("espece_id in ?", ids).Delete(&model.Individu{}).Error; err != nil {
				log.Println(err)
				page.Error = err.Error()
				return nil
			}
		}

		if err := db.Select("Especes").Delete(page.Classification).Error; err != nil {
			log.Println(err)
			page.Error = err.Error()
			return nil
		}
		return e.Close()
	})

	e.OnSubmit("form", func(s *ihui.Session, event ihui.Event) error {
		data := event.Data.(map[string]interface{})
		page.Classification.Nom = data["classification"].(string)
		if page.Classification.Nom == "" {
			page.Error = "Nom obligatoire"
			return nil
		}
		if err := db.Create(page.Classification).Error; err != nil {
			page.Error = err.Error()
			return nil
		}
		return e.Close()
	})

	return e.WriteGoTemplate(TemplatesFs, "templates/classification.html", page)
}
