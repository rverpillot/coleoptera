package pages

import (
	"database/sql"
	"log"
	"strconv"
	"time"

	"github.com/jinzhu/gorm"
	"github.com/rverpillot/coleoptera/model"
	"github.com/rverpillot/ihui"
)

type PageIndividu struct {
	tmpl           *ihui.PageAce
	Individu       model.Individu
	Admin          bool
	Edit           bool
	EditMapCreated bool
	Delete         bool
	Especes        []model.Espece
	Departements   []model.Departement
	Sites          []string
	Communes       []string
	Recolteurs     []string
	Error          string
	Search         string
}

func newPageIndividu(individu model.Individu, editMode bool) *PageIndividu {
	page := &PageIndividu{
		Individu: individu,
		Edit:     editMode,
	}
	page.tmpl = newAceTemplate("individu.ace", page)
	return page
}

func (page *PageIndividu) Render(p ihui.Page) {
	db := p.Get("db").(*gorm.DB)
	page.Especes = model.AllEspeces(db)

	page.Admin = p.Get("admin").(bool)
	page.Sites = model.AllSites(db)
	page.Communes = model.AllCommunes(db)
	page.Departements = model.AllDepartements(db)
	page.Recolteurs = model.AllRecolteurs(db)

	page.tmpl.Render(p)

	p.On("create", "page", func(s *ihui.Session, event ihui.Event) bool {
		if page.Edit {
			s.Script(`createEditMap("#mapedit")`)
		} else {
			s.Script(`createPreviewMap("#mappreview",%f,%f)`, page.Individu.Longitude, page.Individu.Latitude)
		}
		return false
	})

	p.On("update", "page", func(s *ihui.Session, event ihui.Event) bool {
		if page.Edit && !page.EditMapCreated {
			s.Script(`createEditMap("#mapedit")`)
		}
		page.EditMapCreated = true
		return false
	})

	p.On("form", "form", func(s *ihui.Session, event ihui.Event) bool {
		data := event.Data.(map[string]interface{})
		name := data["name"].(string)
		val := data["val"].(string)
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
			s.Script("updateEditMap(%f,%f)", page.Individu.Latitude, page.Individu.Longitude)
		case "latitude":
			lat, _ := strconv.ParseFloat(val, 64)
			page.Individu.Latitude = lat
			page.Search = ""
			s.Script("updateEditMap(%f,%f)", page.Individu.Latitude, page.Individu.Longitude)
		case "recolteur":
			page.Individu.Recolteur = val
		case "commentaire":
			page.Individu.Commentaire = sql.NullString{val, true}
		}
		return true
	})

	p.On("click", "#cancel", func(s *ihui.Session, event ihui.Event) bool {
		return s.CloseModalPage()
	})

	p.On("click", "#edit", func(s *ihui.Session, event ihui.Event) bool {
		page.Edit = true
		return true
	})

	p.On("click", "#delete", func(s *ihui.Session, event ihui.Event) bool {
		page.Delete = true
		return true
	})

	p.On("click", "#confirm-delete", func(s *ihui.Session, event ihui.Event) bool {
		if page.Individu.ID > 0 {
			if err := db.Delete(page.Individu).Error; err != nil {
				log.Println(err)
				page.Error = err.Error()
				return true
			}
		}
		return s.CloseModalPage()
	})

	p.On("click", "#cancel-delete", func(s *ihui.Session, event ihui.Event) bool {
		return s.CloseModalPage()
	})

	p.On("click", "#add-espece", func(s *ihui.Session, event ihui.Event) bool {
		espece := model.Espece{}
		s.ShowPage("espece", newPageEspece(&espece), &ihui.Options{Modal: true})
		if !db.NewRecord(espece) {
			page.Individu.Espece = espece
			page.Individu.EspeceID = espece.ID
		}
		return true
	})

	p.On("click", "#validation", func(s *ihui.Session, event ihui.Event) bool {
		log.Println(page.Individu)
		if page.Individu.Espece.ID == 0 {
			page.Error = "Genre/espèce absent !"
			return true
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
			return true
		}
		return s.CloseModalPage()
	})

	p.On("position", "page", func(s *ihui.Session, event ihui.Event) bool {
		pos := event.Data.(map[string]interface{})
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
		return true
	})

}
