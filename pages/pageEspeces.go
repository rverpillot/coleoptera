package pages

import (
	"github.com/rverpillot/coleoptera/model"
	"github.com/rverpillot/ihui"
	"gorm.io/gorm"
)

type PageEspeces struct {
	menu            *Menu
	Classifications []model.Classification
	Nb              int
	Admin           bool
}

func NewPageEspeces(menu *Menu) *PageEspeces {
	return &PageEspeces{
		menu: menu,
	}
}

func (page *PageEspeces) Render(e *ihui.HTMLElement) error {
	db := e.Get("db").(*gorm.DB)
	page.Admin = e.Get("admin").(bool)
	page.Nb = model.CountAllEspeces(db)
	page.Classifications = model.AllClassifications(db)

	if err := e.WriteGoTemplate(TemplatesFs, "templates/especes.html", page); err != nil {
		return err
	}

	e.On("click", ".espece", func(session *ihui.Session, event ihui.Event) error {
		var espece model.Espece
		db.First(&espece, event.Value())
		session.Set("search_espece", espece.ID)
		return page.menu.ShowPage(session, "individus")
	})

	e.On("click", "#add-espece", func(s *ihui.Session, event ihui.Event) error {
		var espece model.Espece
		return s.ShowModal("espece", newPageEspece(&espece), nil)
	})

	e.On("click", "#add-classification", func(s *ihui.Session, event ihui.Event) error {
		var classification model.Classification
		return s.ShowModal("classification", newPageClassification(&classification), nil)
	})
	return nil

}
