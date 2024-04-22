package pages

import (
	"database/sql"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"path"
	"strconv"
	"strings"
	"time"

	"github.com/jinzhu/gorm"
	"github.com/jung-kurt/gofpdf"
	"github.com/rverpillot/coleoptera/model"
	"github.com/rverpillot/ihui"
)

type PageIndividus struct {
	tmpl          ihui.Template
	menu          *Menu
	selection     map[uint]bool
	SelectCount   int
	AllSelected   bool
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

func (page *PageIndividus) Render(p *ihui.Page) error {
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

	if err := p.WriteAce(TemplatesFs, "templates/individus.ace", page); err != nil {
		return err
	}

	p.On("page-created", "", func(s *ihui.Session, _ ihui.Event) error {
		page.Pagination.SetPage(1)
		return nil
	})

	p.On("input", ".search", func(s *ihui.Session, event ihui.Event) error {
		s.Set("search_individus", event.Value())
		s.Set("search_espece", uint(0))
		page.Pagination.SetPage(1)
		return nil
	})

	p.On("check", ".selectAll", func(s *ihui.Session, event ihui.Event) error {
		page.AllSelected = event.IsChecked()

		if page.AllSelected {
			rows, err := db.Table("individus").Select("id").Rows()
			if err != nil {
				return err
			}
			defer rows.Close()
			for rows.Next() {
				var id uint
				rows.Scan(&id)
				page.selection[id] = true
			}
			page.SelectCount = len(page.selection)
		} else {
			page.clearSelection()
		}
		return nil
	})

	p.On("click", ".detail", func(s *ihui.Session, event ihui.Event) error {
		id := event.Value()
		var individu model.Individu
		db.Preload("Espece").Preload("Departement").Find(&individu, id)
		return s.ShowPage("individu", newPageIndividu(individu, false), &ihui.Options{Modal: true})
	})

	p.On("check", ".select", func(s *ihui.Session, event ihui.Event) error {
		ID, _ := strconv.Atoi(event.Id)
		if event.IsChecked() {
			page.selection[uint(ID)] = true
		} else {
			delete(page.selection, uint(ID))
		}
		page.SelectCount = len(page.selection)
		return nil
	})

	p.On("click", "#reset", func(s *ihui.Session, event ihui.Event) error {
		s.Set("search_individus", "")
		s.Set("search_espece", uint(0))
		page.Pagination.SetPage(1)
		return nil
	})

	p.On("click", "#next", func(s *ihui.Session, event ihui.Event) error {
		page.Pagination.NextPage()
		return nil
	})

	p.On("click", "#previous", func(s *ihui.Session, event ihui.Event) error {
		page.Pagination.PreviousPage()
		return nil
	})

	p.On("click", "table .sortable", func(s *ihui.Session, event ihui.Event) error {
		name := event.Id
		if name == page.fieldSort {
			page.ascendingSort = !page.ascendingSort
		} else {
			page.fieldSort = name
			page.ascendingSort = true
		}
		return nil
	})

	p.On("click", "#add", func(s *ihui.Session, event ihui.Event) error {
		individu := model.Individu{
			Date:      time.Now(),
			Sexe:      "M",
			Latitude:  47.626785,
			Longitude: 6.997305,
			Altitude:  sql.NullInt64{Int64: 100, Valid: true},
		}
		return s.ShowPage("individu", newPageIndividu(individu, true), &ihui.Options{Modal: true})
	})

	p.On("click", "#printLabels", func(s *ihui.Session, event ihui.Event) error {
		tmpDir := path.Join(os.TempDir(), "coleoptera")

		f, err := ioutil.TempFile(tmpDir, "etiquettes-*.pdf")
		if err != nil {
			return err
		}
		defer f.Close()

		if err := page.printLabels(db, f); err != nil {
			return err
		}
		s.Script(`
		win = window.open("","print")
		if (win) {win.location = "tmp/%s"}
		`, path.Base(f.Name()))

		page.clearSelection()
		return nil
	})

	p.On("click", "#export", func(s *ihui.Session, event ihui.Event) error {
		tmpDir := path.Join(os.TempDir(), "coleoptera")

		f, err := ioutil.TempFile(tmpDir, "coleoptera-*.csv")
		if err != nil {
			return err
		}
		defer f.Close()

		if err := export(db, f); err != nil {
			return err
		}
		return s.Script(`window.open("tmp/%s","export")`, path.Base(f.Name()))
	})

	return nil
}

func export(db *gorm.DB, output io.Writer) error {
	var count int
	var individus []model.Individu

	if err := db.Model(&model.Individu{}).Count(&count).
		Preload("Espece.Classification").
		Preload("Espece").
		Find(&individus).Error; err != nil {
		return err
	}
	headers := []string{
		"Classification",
		"Ordre",
		"Espece",
		"Site",
		"GPS",
		"Altitude",
		"Commune",
		"Code",
		"Sexe",
		"Date",
		"Commentaire",
		"Recolteur",
	}
	output.Write([]byte(strings.Join(headers, "\t") + "\n"))

	for _, individu := range individus {
		ordre := ""
		if individu.Espece.Classification.Ordre.Valid {
			ordre = fmt.Sprintf("%d", individu.Espece.Classification.Ordre.Int64)
		}
		data := []string{
			individu.Espece.Classification.Nom,
			ordre,
			individu.Espece.NomEspece(),
			individu.Site,
			individu.Localization(),
			fmt.Sprintf("%d", individu.Altitude.Int64),
			individu.Commune,
			individu.Code,
			individu.Sexe,
			individu.Date.Format("02/01/2006"),
			strings.Replace(individu.Commentaire.String, "\n", "", -1),
			individu.Recolteur,
		}
		output.Write([]byte(strings.Join(data, "\t") + "\n"))
	}
	return nil
}

func (page *PageIndividus) clearSelection() {
	page.selection = make(map[uint]bool)
	page.SelectCount = 0
	page.AllSelected = false
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
			pdf.Circle(x+1, y+height/2, 0.2, "F")
			pdf.SetXY(x, y)
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

func printText(pdf *gofpdf.Fpdf, width, height float64, text string) {
	tr := pdf.UnicodeTranslatorFromDescriptor("")
	text = tr(text)
	for size := float64(5); size > 4; size -= 0.5 {
		pdf.SetFontSize(size)
		if pdf.GetStringWidth(text) < width {
			break
		}
	}
	pdf.CellFormat(width, height, text, "", 2, "C", false, 0, "")
}

func printLabel1(pdf *gofpdf.Fpdf, x, y, width, height float64, individu *model.Individu) {
	printText(pdf, width, 2.5, individu.Commune)
	printText(pdf, width, 2.5, "("+individu.Departement.Nom+")")
	printText(pdf, width, 2.5, "France")
	printText(pdf, width, 2.5, individu.Recolteur)
}

func printLabel2(pdf *gofpdf.Fpdf, x, y, width, height float64, individu *model.Individu) {
	printText(pdf, width, 3.3, individu.Site)
	printText(pdf, width, 3.3, fmt.Sprintf("%dm", individu.Altitude.Int64))
	printText(pdf, width, 3.3, fmtDate(individu.Date))
}

func fmtDate(date time.Time) string {
	romans := []string{

		"", "I", "II", "III", "IV", "V", "VI", "VII", "VIII", "IX", "X", "XI", "XII",
	}
	return fmt.Sprintf("%d-%s-%d", date.Day(), romans[date.Month()], date.Year())
}
