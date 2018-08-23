package pages

import (
	"github.com/jinzhu/gorm"
	"rverpi/coleoptera.v3/model"
	"rverpi/ihui.v2"
)

type PageEspeces struct {
	tmpl            *ihui.PageAce
	menu            *Menu
	Classifications []model.Classification
}

func NewPageEspeces(menu *Menu) *PageEspeces {
	page := &PageEspeces{
		menu: menu,
	}
	page.tmpl = newAceTemplate("especes.ace", page)
	return page
}

func (page *PageEspeces) Render(p ihui.Page) {
	db := p.Get("db").(*gorm.DB)
	page.Classifications = model.AllClassifications(db)

	page.tmpl.Render(p)
	p.Add("#menu", page.menu)

	p.On("click", ".espece", func(session *ihui.Session, event ihui.Event) {
		var espece model.Espece
		db.First(&espece, event.Value())
		session.Set("search_espece", espece.ID)
	})
}
