package pages

import (
	"database/sql"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"path"
	"strconv"
	"time"

	"bitbucket.org/rverpi90/coleoptera.v3/model"
	"bitbucket.org/rverpi90/ihui"
	"github.com/jinzhu/gorm"
	"github.com/jung-kurt/gofpdf"
)

type PageIndividus struct {
	tmpl          *ihui.PageAce
	menu          *Menu
	selection     map[uint]bool
	SelectCount   int
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
		page.SelectCount = len(page.selection)
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

	p.On("click", "#printLabels", func(s *ihui.Session, event ihui.Event) bool {
		f, err := ioutil.TempFile("", "coleoptera*.pdf")
		if err != nil {
			log.Print(err)
			return false
		}
		defer f.Close()

		if err := page.printLabels(db, f); err != nil {
			log.Print(err)
			return false
		}
		s.Script(`
		win = window.open("","print")
		if (win) {win.location = "/pdf/%s"}
		`, path.Base(f.Name()))
		return false
	})
}

func (page *PageIndividus) printLabels(db *gorm.DB, output io.Writer) error {
	const width = 20
	const height = 10
	const cols = 8
	const rows = 25
	const leftMargin = (210 - width*cols) / 2
	const topMargin = (297 - height*rows) / 2

	pdf := gofpdf.New("Portrait", "mm", "A4", "")
	defer pdf.Close()
	pdf.SetFont("Helvetica", "", 5)
	pdf.SetLineWidth(0.1)

	col := 0
	row := 0
	for id := range page.selection {
		var individu model.Individu
		if err := db.Preload("Departement").First(&individu, id).Error; err != nil {
			continue
		}

		printLabels := []func(*gofpdf.Fpdf, float64, float64, float64, float64, *model.Individu){
			printLabel1,
			printLabel2,
		}

		for _, printLabel := range printLabels {
			if col == 0 && row == 0 {
				pdf.AddPage()
			}

			x := float64(leftMargin + col*width)
			y := float64(topMargin + row*height)
			pdf.Rect(x, y, width, height, "D")
			printLabel(pdf, x, y, width, height, &individu)

			col++
			if col >= cols {
				col = 0
				row++
			}
			if row >= rows {
				row = 0
			}
		}
	}
	return pdf.Output(output)
}

func printLabel1(pdf *gofpdf.Fpdf, x, y, width, height float64, individu *model.Individu) {
	tr := pdf.UnicodeTranslatorFromDescriptor("")
	pdf.Circle(x+2, y+height/2, 0.2, "F")
	pdf.SetXY(x, y)
	pdf.CellFormat(width, 2.5, tr(individu.Commune), "", 2, "C", false, 0, "")
	pdf.CellFormat(width, 2.5, tr(individu.Departement.Nom), "", 2, "C", false, 0, "")
	pdf.CellFormat(width, 2.5, "France", "", 2, "C", false, 0, "")
	pdf.CellFormat(width, 2.5, tr(individu.Recolteur), "", 2, "C", false, 0, "")
}

func printLabel2(pdf *gofpdf.Fpdf, x, y, width, height float64, individu *model.Individu) {
	tr := pdf.UnicodeTranslatorFromDescriptor("")
	pdf.Circle(x+2, y+5, 0.2, "F")
	pdf.SetXY(x, y)
	pdf.CellFormat(width, 3.3, tr(individu.Site), "", 2, "C", false, 0, "")
	pdf.CellFormat(width, 3.3, fmt.Sprintf("%dm", individu.Altitude.Int64), "", 2, "C", false, 0, "")
	pdf.CellFormat(width, 3.3, fmtDate(individu.Date), "", 2, "C", false, 0, "")
}

func fmtDate(date time.Time) string {
	romans := []string{
		"I", "II", "III", "IV", "V", "VI", "VII", "VIII", "IX", "X", "XI", "XII",
	}
	return fmt.Sprintf("%d-%s-%d", date.Day(), romans[date.Month()], date.Year())
}
