package pages

import (
	"database/sql"
	"strconv"
	"time"

	"github.com/jinzhu/gorm"
	"rverpi90/coleoptera.v3/model"
	"bitbucket.org/rverpi90/ihui"
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
	fieldSort     string
	ascendingSort bool
}

func NewPageIndividus(menu *Menu) *PageIndividus {
	return &PageIndividus{
		tmpl:       newAceTemplate("individus.ace", nil),
		menu:       menu,
		selection:  make(map[uint]bool),
		Pagination: ihui.NewPaginator(60),
		fieldSort:  "date",
	}
}

func (page *PageIndividus) ShowSort(name string) string {
	if name == page.fieldSort {
		if page.ascendingSort {
			return "sortable sorted ascending"
		} else {
			return "sortable sorted descending"
		}
	}
	return "sortable"
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

	order := page.fieldSort
	if !page.ascendingSort {
		order += " desc"
	}
	total := model.LoadIndividus(db, &page.Individus, page.Pagination.Current.Index, page.Pagination.PageSize, page.Search, espece_id, order)

	page.Pagination.SetTotal(total)

	for n, i := range page.Individus {
		if page.selection[i.ID] {
			page.Individus[n].Selected = true
		}
	}

	page.tmpl.SetModel(page)
	page.tmpl.Render(p)
	p.Add("#menu", page.menu)

	p.On("create", "page", func(s *ihui.Session, _ ihui.Event) bool {
		page.Pagination.SetPage(1)
		return false
	})

	p.On("input", ".search", func(s *ihui.Session, event ihui.Event) bool {
		s.Set("search_individus", event.Value())
		s.Set("search_espece", uint(0))
		page.Pagination.SetPage(1)
		return true
	})

	p.On("click", ".detail", func(s *ihui.Session, event ihui.Event) bool {
		id := event.Value()
		var individu model.Individu
		db.Preload("Espece").Preload("Departement").Find(&individu, id)
		return s.ShowPage("individu", newPageIndividu(individu, false), &ihui.Options{Modal: true})
	})

	p.On("check", ".select", func(s *ihui.Session, event ihui.Event) bool {
		ID, _ := strconv.Atoi(event.Id)
		if event.Data.(bool) {
			page.selection[uint(ID)] = true
		} else {
			delete(page.selection, uint(ID))
		}
		return true
	})

	p.On("click", "#reset", func(s *ihui.Session, event ihui.Event) bool {
		s.Set("search_individus", "")
		s.Set("search_espece", uint(0))
		page.Pagination.SetPage(1)
		return true
	})

	p.On("click", "#next", func(s *ihui.Session, event ihui.Event) bool {
		page.Pagination.NextPage()
		return true
	})

	p.On("click", "#previous", func(s *ihui.Session, event ihui.Event) bool {
		page.Pagination.PreviousPage()
		return true
	})

	p.On("click", "table .sortable", func(s *ihui.Session, event ihui.Event) bool {
		name := event.Id
		if name == page.fieldSort {
			page.ascendingSort = !page.ascendingSort
		} else {
			page.fieldSort = name
			page.ascendingSort = true
		}
		return true
	})

	p.On("click", "#add", func(s *ihui.Session, event ihui.Event) bool {
		individu := model.Individu{
			Date:      time.Now(),
			Sexe:      "M",
			Latitude:  47.626785,
			Longitude: 6.997305,
			Altitude:  sql.NullInt64{100, true},
		}
		return s.ShowPage("individu", newPageIndividu(individu, true), &ihui.Options{Modal: true})
	})
}
