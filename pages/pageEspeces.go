package pages

import (
	"github.com/jinzhu/gorm"
	"github.com/rverpillot/coleoptera/model"
	"github.com/rverpillot/ihui"
)

type PageEspeces struct {
	tmpl            *ihui.PageAce
	menu            *Menu
	Classifications []model.Classification
	Nb              int
}

func NewPageEspeces(menu *Menu) *PageEspeces {
	page := &PageEspeces{menu: menu}
	page.tmpl = newAceTemplate("especes.ace", page)
	return page
}

func (page *PageEspeces) Render(p *ihui.Page) {
	db := p.Get("db").(*gorm.DB)
	page.Nb = model.CountAllEspeces(db)
	page.Classifications = model.AllClassifications(db)

	page.tmpl.Render(p)

	p.On("click", ".espece", func(session *ihui.Session, event ihui.Event) {
		var espece model.Espece
		db.First(&espece, event.Value())
		session.Set("search_espece", espece.ID)
		page.menu.ShowPage(session, "individus")
	})
}
