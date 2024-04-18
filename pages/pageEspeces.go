package pages

import (
	"github.com/jinzhu/gorm"
	"github.com/rverpillot/coleoptera/model"
	"github.com/rverpillot/ihui"
	"github.com/rverpillot/ihui/templating"
)

type PageEspeces struct {
	tmpl            *templating.PageAce
	menu            *Menu
	Classifications []model.Classification
	Nb              int
}

func NewPageEspeces(menu *Menu) *PageEspeces {
	page := &PageEspeces{menu: menu}
	page.tmpl = newAceTemplate("especes.ace", page)
	return page
}

func (page *PageEspeces) Render(p *ihui.Page) error {
	db := p.Get("db").(*gorm.DB)
	page.Nb = model.CountAllEspeces(db)
	page.Classifications = model.AllClassifications(db)

	if err := page.tmpl.Render(p); err != nil {
		return err
	}

	p.On("click", ".espece", func(session *ihui.Session, event ihui.Event) error {
		var espece model.Espece
		db.First(&espece, event.Value())
		session.Set("search_espece", espece.ID)
		return page.menu.ShowPage(session, "individus")
	})

	return nil
}
