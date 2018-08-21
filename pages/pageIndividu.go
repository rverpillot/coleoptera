package pages

import (
	"database/sql"
	"io"
	"log"
	"rverpi/coleoptera/model"
	"rverpi/ihui"
	"strconv"
	"time"

	"github.com/jinzhu/gorm"
)

type PageIndividu struct {
	*Page
	Individu            model.Individu
	Admin               bool
	Edit                bool
	Delete              bool
	Especes             []model.Espece
	Departements        []model.Departement
	Sites               []string
	Communes            []string
	Recolteurs          []string
	Error               string
	Search              string
	SearchAction        ihui.ChangeAction
	ValidAction         ihui.ClickAction
	AddEspeceAction     ihui.ClickAction
	CloseAction         ihui.ClickAction
	EditAction          ihui.ClickAction
	DeleteAction        ihui.ClickAction
	ConfirmDeleteAction ihui.ClickAction
	CancelDeleteAction  ihui.ClickAction
	FormAction          ihui.FormAction
}

func newPageIndividu(individu model.Individu, editMode bool) ihui.PageRender {
	page := &PageIndividu{
		Page:     NewPage("individu", true),
		Individu: individu,
		Edit:     editMode,
	}
	return page
}

func (page *PageIndividu) OnInit(ctx *ihui.Context) {
	db := ctx.Get("db").(*gorm.DB)

	page.Admin = ctx.Get("admin").(bool)
	page.Sites = model.AllSites(db)
	page.Communes = model.AllCommunes(db)
	page.Departements = model.AllDepartements(db)
	page.Recolteurs = model.AllRecolteurs(db)

	page.FormAction = func(name string, val string) {
		log.Println(name, val)
		switch name {
		case "date":
			page.Individu.Date, _ = time.Parse("02/01/2006", val)
		case "espece":
			id, _ := strconv.Atoi(val)
			page.Individu.EspeceID = uint(id)
			db.First(&page.Individu.Espece, id)
		case "sexe":
			page.Individu.Sexe = val
		case "site":
			page.Individu.Site = val
		case "altitude":
			altitude, _ := strconv.Atoi(val)
			page.Individu.Altitude = sql.NullInt64{Int64: int64(altitude), Valid: true}
		case "commune":
			page.Individu.Commune = val
			lat, lng, err := model.FindLatLng(page.Individu.Commune)
			if err != nil {
				log.Println(err)
				break
			}
			page.Individu.Latitude = lat
			page.Individu.Longitude = lng
		case "longitude":
			long, _ := strconv.ParseFloat(val, 64)
			page.Individu.Longitude = long
			page.Search = ""
			ctx.Script("updateEditMap(%f,%f)", page.Individu.Latitude, page.Individu.Longitude)
		case "latitude":
			lat, _ := strconv.ParseFloat(val, 64)
			page.Individu.Latitude = lat
			page.Search = ""
			ctx.Script("updateEditMap(%f,%f)", page.Individu.Latitude, page.Individu.Longitude)
		case "recolteur":
			page.Individu.Recolteur = val
		case "commentaire":
			page.Individu.Commentaire = sql.NullString{val, true}
		}
	}

	page.SearchAction = func(val string) {
		log.Println(val)
		page.Individu.Latitude, page.Individu.Longitude, _ = model.FindLatLng(val)
		ctx.Script("updateEditMap(%f,%f)", page.Individu.Latitude, page.Individu.Longitude)
	}

	page.ValidAction = func(_ string) {
		log.Println(page.Individu)
		if page.Individu.Espece.ID == 0 {
			page.Error = "Genre/espèce absent !"
			return
		}
		var err error
		if db.NewRecord(page.Individu) {
			err = db.Set("gorm:save_associations", false).Create(&page.Individu).Error
		} else {
			err = db.Set("gorm:save_associations", false).Save(&page.Individu).Error
		}
		if err != nil {
			log.Println(err)
			page.Error = err.Error()
			return
		}
		page.Close()
	}

	page.AddEspeceAction = func(_ string) {
		espece := model.Espece{}
		ctx.DisplayPage(newPageEspece(&espece), true)
		if !db.NewRecord(espece) {
			page.Individu.Espece = espece
			page.Individu.EspeceID = espece.ID
		}
	}

	page.CloseAction = func(_ string) {
		page.Close()
	}
	page.EditAction = func(_ string) {
		page.Edit = true
	}
	page.DeleteAction = func(_ string) {
		page.Delete = true
	}
	page.CancelDeleteAction = func(_ string) {
		page.Delete = false
	}
	page.ConfirmDeleteAction = func(_ string) {
		if page.Individu.ID > 0 {
			if err := db.Delete(page.Individu).Error; err != nil {
				log.Println(err)
				page.Error = err.Error()
				return
			}
		}
		page.Close()
	}

	page.On("position", func(ctx *ihui.Context) {
		pos := ctx.Event.Data.(map[string]interface{})
		log.Println(pos)
		page.Individu.Latitude = pos["lat"].(float64)
		page.Individu.Longitude = pos["lng"].(float64)

		var altitude int64
		var err error
		page.Individu.Commune, page.Individu.Code, altitude, err = model.FindLocation(page.Individu.Latitude, page.Individu.Longitude)
		if err != nil {
			log.Println(err)
		}
		page.Individu.Altitude = sql.NullInt64{altitude, true}
	})

}

func (page *PageIndividu) Render(w io.Writer, ctx *ihui.Context) {
	db := ctx.Get("db").(*gorm.DB)
	page.Especes = model.AllEspeces(db)

	page.renderPage(w, page)
}
