package pages

import (
	"log"
	"strconv"
	"strings"

	"github.com/rverpillot/coleoptera/model"
	"github.com/rverpillot/ihui"

	"github.com/jinzhu/gorm"
)

type PageEspece struct {
	tmpl            ihui.Template
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
		tmpl:   newAceTemplate("espece.ace"),
		Espece: espece,
	}
}

func (page *PageEspece) ID() string {
	if page.Espece.ID == 0 {
		return ""
	}
	return strconv.Itoa(int(page.Espece.ID))
}

func (page *PageEspece) Render(p *ihui.Page) error {
	db := p.Get("db").(*gorm.DB)

	page.Classifications = model.AllClassifications(db)
	page.AllGenres = model.AllGenres(db)
	page.AllSousGenres = model.AllSousGenres(db)
	page.AllEspeces = model.AllNomEspeces(db)
	page.AllSousEspeces = model.AllSousEspeces(db)

	if err := p.WriteTemplate(page.tmpl, page); err != nil {
		return err
	}

	p.On("click", "[id=add-classification]", func(s *ihui.Session, ev ihui.Event) error {
		var classification model.Classification
		s.ShowPage("classification", newPageClassification(&classification), &ihui.Options{Modal: true})
		if !db.NewRecord(classification) {
			page.Espece.Classification = classification
			page.Espece.ClassificationID = classification.ID
		}
		return nil
	})

	p.On("click", "[id=cancel]", func(s *ihui.Session, ev ihui.Event) error {
		return p.Close()
	})

	p.On("submit", "form", func(s *ihui.Session, ev ihui.Event) error {
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
		}
		return p.Close()
	})

	return nil
}
