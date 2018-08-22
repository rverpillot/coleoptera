package pages

import (
	"database/sql"
	"strconv"
	"time"

	"github.com/jinzhu/gorm"
	"rverpi/coleoptera.v3/model"
	"rverpi/ihui.v2"
)

type PageIndividus struct {
	tmpl          *ihui.AceTemplateDrawer
	menu          *Menu
	selection     map[uint]bool
	Pagination    *ihui.Paginator
	Individus     []model.Individu
	Admin         bool
	Search        string
	ShowAllButton bool
}

func NewPageIndividus(menu *Menu) *PageIndividus {
	page := &PageIndividus{
		menu:       menu,
		selection:  make(map[uint]bool),
		Pagination: ihui.NewPaginator(60),
	}
	page.tmpl = newAceTemplate("individus.ace", page)
	return page
}

func (page *PageIndividus) Draw(p ihui.Page) {
	db := p.Get("db").(*gorm.DB)

	var espece_id uint

	page.Admin = ctx.Get("admin").(bool)

	if p.Get("search_individus") != nil {
		page.Search = ctx.Get("search_individus").(string)
	}
	if p.Get("search_espece") != nil {
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

	p.Draw(page.tmpl)

	p.On("load", "page", func(s *ihui.Session, _ ihui.Event) {
		page.Pagination.SetPage(1)
	})

	p.On("input", ".search", func(s *ihui.Session, event ihui.Event) {
		s.Set("search_individus", event.Value())
		s.Set("search_espece", uint(0))
		s.Pagination.SetPage(1)
	})

	p.On("click", ".detail", func(s *ihui.Session, event ihui.Event) {
		id := event.Value()
		var individu model.Individu
		db.Preload("Espece").Preload("Departement").Find(&individu, id)
		s.ShowPage(newPageIndividu(individu, false), ihui.Options{Modal: true})
	})

	p.On("change", ".select", func(s *ihui.Session, event ihui.Event) {
		ID, _ := strconv.Atoi(event.Source())
		if val {
			page.selection[uint(ID)] = true
		} else {
			delete(page.selection, uint(ID))
		}
	})

	p.On("click", "[id=reset]", func(s *ihui.Session, event ihui.Event) {
		s.Set("search_individus", "")
		s.Set("search_espece", uint(0))
		page.Pagination.SetPage(1)
	})

	p.On("click", "[id=next]", func(s *ihui.Session, event ihui.Event) {
		page.Pagination.NextPage()
	})

	p.On("click", "[id=previous]", func(s *ihui.Session, event ihui.Event) {
		page.Pagination.PreviousPage()
	})

	p.On("click", "[id=add]", func(s *ihui.Session, event ihui.Event) {
		individu := model.Individu{
			Date:      time.Now(),
			Sexe:      "M",
			Latitude:  47.626785,
			Longitude: 6.997305,
			Altitude:  sql.NullInt64{100, true},
		}
		s.ShowPage(newPageIndividu(individu, true), ihui.Options{Modal: true})
	})
}
