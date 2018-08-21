package pages

import (
	"database/sql"
	"io"
	"rverpi/coleoptera/model"
	"rverpi/ihui"
	"strconv"
	"time"

	"github.com/jinzhu/gorm"
)

type PageIndividus struct {
	ihui.Page
	menu               *Menu
	selection          map[uint]bool
	Pagination         *ihui.Paginator
	Individus          []model.Individu
	Admin              bool
	Search             string
	ShowAllButton      bool
	SearchAction       ihui.ChangeAction
	SelectAction       ihui.CheckAction
	ResetAction        ihui.ClickAction
	DetailAction       ihui.ClickAction
	PreviousPageAction ihui.ClickAction
	NextPageAction     ihui.ClickAction
	AddAction          ihui.ClickAction
}

func NewPageIndividus(menu *Menu) *PageIndividus {
	page := &PageIndividus{
		menu:       menu,
		selection:  make(map[uint]bool),
		Pagination: ihui.NewPaginator(60),
	}
	page.Add("#menu", menu)
	return page
}

func (page *PageIndividus) OnInit(ctx *ihui.Context) {
	db := ctx.Get("db").(*gorm.DB)

	page.DetailAction = func(id string) {
		var individu model.Individu
		db.Preload("Espece").Preload("Departement").Find(&individu, id)
		ctx.DisplayPage(newPageIndividu(individu, false), true)
	}

	page.AddAction = func(_ string) {
		individu := model.Individu{
			Date:      time.Now(),
			Sexe:      "M",
			Latitude:  47.626785,
			Longitude: 6.997305,
			Altitude:  sql.NullInt64{100, true},
		}
		ctx.DisplayPage(newPageIndividu(individu, true), true)
	}

	page.SelectAction = func(val bool, id string) {
		ID, _ := strconv.Atoi(id)
		if val {
			page.selection[uint(ID)] = true
		} else {
			delete(page.selection, uint(ID))
		}
	}

	page.SearchAction = func(val string) {
		ctx.Set("search_individus", val)
		ctx.Set("search_espece", uint(0))
		page.Pagination.SetPage(1)
		page.Trigger("searching", ctx)
	}

	page.ResetAction = func(_ string) {
		ctx.Set("search_individus", "")
		ctx.Set("search_espece", uint(0))
		page.Pagination.SetPage(1)
		page.Trigger("searching", ctx)
	}

	page.PreviousPageAction = func(_ string) {
		page.Pagination.PreviousPage()
	}
	page.NextPageAction = func(_ string) {
		page.Pagination.NextPage()
	}

}

func (page *PageIndividus) OnShow(ctx *ihui.Context) {
	page.Pagination.SetPage(1)
}

func (page *PageIndividus) Render(w io.Writer, ctx *ihui.Context) {
	db := ctx.Get("db").(*gorm.DB)
	var espece_id uint

	page.Admin = ctx.Get("admin").(bool)

	if ctx.Get("search_individus") != nil {
		page.Search = ctx.Get("search_individus").(string)
	}
	if ctx.Get("search_espece") != nil {
		espece_id = ctx.Get("search_espece").(uint)
		if espece_id != 0 {
			page.Search = ""
		}
	}

	page.ShowAllButton = page.Search != "" || espece_id != 0

	total := model.LoadIndividus(db, &page.Individus, page.Pagination.Current.Index, page.Pagination.PageSize, page.Search, espece_id)

	page.Pagination.SetTotal(total)

	for n, i := range page.Individus {
		if page.selection[i.ID] {
			page.Individus[n].Selected = true
		}
	}

	renderTemplate("individus", w, page)
}
