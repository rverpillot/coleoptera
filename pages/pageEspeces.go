package pages

import (
	"github.com/jinzhu/gorm"
	"rverpi/coleoptera.v3/model"
	"rverpi/ihui.v2"
)

type PageEspeces struct {
	tmpl            *ihui.AceTemplateDrawer
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

func (page *PageEspeces) Draw(p ihui.PageDrawer) {
	db := p.Get("db").(*gorm.DB)
	page.Classifications = model.AllClassifications(db)

	p.Draw(page.tmpl)

	p.On("click", ".espece", func(session *ihui.Session) {
		var espece model.Espece
		db.First(&espece, id)
		session.Set("search_espece", espece.ID)
	})
}
