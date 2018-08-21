package pages

import (
	"io"
	"rverpi/coleoptera/model"
	"rverpi/ihui"

	"github.com/jinzhu/gorm"
)

type PageEspeces struct {
	*Page
	menu            *Menu
	Classifications []model.Classification
	SelectAction    ihui.ClickAction
}

func NewPageEspeces(menu *Menu) *PageEspeces {
	page := &PageEspeces{
		Page: NewPage("especes", false),
		menu: menu,
	}
	page.Add("#menu", menu)
	return page
}

func (page *PageEspeces) OnInit(ctx *ihui.Context) {
	db := ctx.Get("db").(*gorm.DB)
	page.SelectAction = func(id string) {
		var espece model.Espece
		db.First(&espece, id)
		ctx.Event.Name = "search_espece"
		ctx.Event.Data = espece.ID
		page.Trigger("search_espece", ctx)
	}
}

func (page *PageEspeces) Render(w io.Writer, ctx *ihui.Context) {
	db := ctx.Get("db").(*gorm.DB)
	page.Classifications = model.AllClassifications(db)

	page.renderPage(w, page)

}
