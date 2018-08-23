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
	tmpl          *ihui.PageAce
	menu          *Menu
	selection     map[uint]bool
	Pagination    *ihui.Paginator
	Individus     []model.Individu
	Admin         bool
	Search        string
	ShowAllButton bool
}

func NewPageIndividus(menu *Menu) *PageIndividus {
	return &PageIndividus{
		tmpl:       newAceTemplate("individus.ace", nil),
		menu:       menu,
		selection:  make(map[uint]bool),
		Pagination: ihui.NewPaginator(60),
	}
}

func (page *PageIndividus) Render(p ihui.Page) {
	db := p.Get("db").(*gorm.DB)

	var espece_id uint

	page.Admin = p.Get("admin").(bool)

	if p.Get("search_individus") != nil {
		page.Search = p.Get("search_individus").(string)
	}
	if p.Get("search_espece") != nil {
		espece_id = p.Get("search_espece").(uint)
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

	page.tmpl.SetModel(page)
	page.tmpl.Render(p)
	p.Add("#menu", page.menu)

	p.On("load", "page", func(s *ihui.Session, _ ihui.Event) {
		page.Pagination.SetPage(1)
	})

	p.On("input", ".search", func(s *ihui.Session, event ihui.Event) {
		s.Set("search_individus", event.Value())
		s.Set("search_espece", uint(0))
		page.Pagination.SetPage(1)
	})

	p.On("click", ".detail", func(s *ihui.Session, event ihui.Event) {
		id := event.Value()
		var individu model.Individu
		db.Preload("Espece").Preload("Departement").Find(&individu, id)
		s.ShowPage(newPageIndividu(individu, false), &ihui.Options{Modal: true})
	})

	p.On("change", ".select", func(s *ihui.Session, event ihui.Event) {
		ID, _ := strconv.Atoi(event.Source)
		if event.Data.(bool) {
			page.selection[uint(ID)] = true
		} else {
			delete(page.selection, uint(ID))
		}
	})

	p.On("click", "#reset", func(s *ihui.Session, event ihui.Event) {
		s.Set("search_individus", "")
		s.Set("search_espece", uint(0))
		page.Pagination.SetPage(1)
	})

	p.On("click", "#next", func(s *ihui.Session, event ihui.Event) {
		page.Pagination.NextPage()
	})

	p.On("click", "#previous", func(s *ihui.Session, event ihui.Event) {
		page.Pagination.PreviousPage()
	})

	p.On("click", "#add", func(s *ihui.Session, event ihui.Event) {
		individu := model.Individu{
			Date:      time.Now(),
			Sexe:      "M",
			Latitude:  47.626785,
			Longitude: 6.997305,
			Altitude:  sql.NullInt64{100, true},
		}
		s.ShowPage(newPageIndividu(individu, true), &ihui.Options{Modal: true})
	})
}
