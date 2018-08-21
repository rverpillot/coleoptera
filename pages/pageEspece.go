package pages

import (
	"io"
	"log"
	"rverpi/coleoptera/model"
	"rverpi/ihui"
	"strconv"
	"strings"

	"github.com/jinzhu/gorm"
)

type PageEspece struct {
	*Page
	Espece                  *model.Espece
	Classifications         []model.Classification
	AllGenres               []string
	AllSousGenres           []string
	AllEspeces              []string
	AllSousEspeces          []string
	Error                   string
	AddClassificationAction ihui.ClickAction
	SubmitAction            ihui.SubmitAction
	CloseAction             ihui.ClickAction
}

func newPageEspece(espece *model.Espece) ihui.PageRender {
	page := &PageEspece{
		Page:   NewPage("espece", true),
		Espece: espece,
	}
	return page
}

func (page *PageEspece) ID() string {
	if page.Espece.ID == 0 {
		return ""
	}
	return strconv.Itoa(int(page.Espece.ID))
}

func (page *PageEspece) OnInit(ctx *ihui.Context) {
	db := ctx.Get("db").(*gorm.DB)

	page.AllGenres = model.AllGenres(db)
	page.AllSousGenres = model.AllSousGenres(db)
	page.AllEspeces = model.AllNomEspeces(db)
	page.AllSousEspeces = model.AllSousEspeces(db)

	page.AddClassificationAction = func(_ string) {
		var classification model.Classification
		ctx.DisplayPage(newPageClassification(&classification), true)
		if !db.NewRecord(classification) {
			page.Espece.Classification = classification
			page.Espece.ClassificationID = classification.ID
		}
	}

	page.CloseAction = func(_ string) {
		page.Close()
	}

	page.SubmitAction = func(data map[string]interface{}) {
		id, _ := strconv.Atoi(data["classification"].(string))

		page.Espece.ClassificationID = uint(id)
		page.Espece.Genre = data["genre"].(string)
		page.Espece.SousGenre = strings.Trim(data["sous_genre"].(string), "()")

		page.Espece.Espece = data["espece"].(string)
		page.Espece.SousEspece = data["sous_espece"].(string)
		page.Espece.Descripteur = data["descripteur"].(string)

		if page.Espece.ClassificationID == 0 || page.Espece.Genre == "" || page.Espece.Espece == "" || page.Espece.Descripteur == "" {
			page.Error = "Informations incomplètes !"
			return
		}

		log.Println(page.Espece)
		if err := db.Create(page.Espece).Error; err != nil {
			log.Println(err)
			page.Error = err.Error()
			return
		}
		page.Close()
	}

}

func (page *PageEspece) Render(w io.Writer, ctx *ihui.Context) {
	db := ctx.Get("db").(*gorm.DB)

	page.Classifications = model.AllClassifications(db)
	page.renderPage(w, page)

}
