package pages

import (
	"log"
	"strconv"
	"strings"

	"github.com/rverpillot/coleoptera/model"
	"github.com/rverpillot/ihui"

	"gorm.io/gorm"
)

type PageEspece struct {
	Espece          *model.Espece
	Classifications []model.Classification
	AllGenres       []string
	AllSousGenres   []string
	AllEspeces      []string
	AllSousEspeces  []string
	Error           string
}

func newPageEspece(espece *model.Espece) *PageEspece {
	return &PageEspece{
		Espece: espece,
	}
}

func (page *PageEspece) ID() string {
	if page.Espece.ID == 0 {
		return ""
	}
	return strconv.Itoa(int(page.Espece.ID))
}

func (page *PageEspece) Render(e *ihui.HTMLElement) error {
	db := e.Get("db").(*gorm.DB)

	page.Classifications = model.AllClassifications(db)
	page.AllGenres = model.AllGenres(db)
	page.AllSousGenres = model.AllSousGenres(db)
	page.AllEspeces = model.AllNomEspeces(db)
	page.AllSousEspeces = model.AllSousEspeces(db)

	if err := e.WriteGoTemplate(TemplatesFs, "templates/espece.html", page); err != nil {
		return err
	}

	e.On("click", "#cancel", func(s *ihui.Session, ev ihui.Event) error {
		return e.Close()
	})

	e.On("submit", "form", func(s *ihui.Session, ev ihui.Event) error {
		data := ev.Data.(map[string]interface{})
		id, _ := strconv.Atoi(data["classification"].(string))

		page.Espece.ClassificationID = uint(id)
		page.Espece.Genre = strings.Title(data["genre"].(string))
		page.Espece.SousGenre = strings.Title(strings.Trim(data["sous_genre"].(string), "()"))

		page.Espece.Espece = data["espece"].(string)
		page.Espece.SousEspece = data["sous_espece"].(string)
		page.Espece.Descripteur = data["descripteur"].(string)

		if page.Espece.ClassificationID == 0 || page.Espece.Genre == "" || page.Espece.Espece == "" || page.Espece.Descripteur == "" {
			page.Error = "Informations incomplètes !"
		}

		log.Println(page.Espece)
		if err := db.Create(page.Espece).Error; err != nil {
			log.Println(err)
			page.Error = err.Error()
			return nil
		}
		return e.Close()
	})

	return nil
}
