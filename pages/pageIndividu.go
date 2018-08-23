package pages

import (
	"database/sql"
	"log"
	"strconv"
	"time"

	"github.com/jinzhu/gorm"
	"rverpi/coleoptera.v3/model"
	"rverpi/ihui.v2"
)

type PageIndividu struct {
	tmpl         *ihui.PageAce
	Individu     model.Individu
	Admin        bool
	Edit         bool
	Delete       bool
	Especes      []model.Espece
	Departements []model.Departement
	Sites        []string
	Communes     []string
	Recolteurs   []string
	Error        string
	Search       string
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

	page.tmpl.Render(p)

	p.On("load", "page", func(s *ihui.Session, event ihui.Event) {
		page.Admin = s.Get("admin").(bool)
		page.Sites = model.AllSites(db)
		page.Communes = model.AllCommunes(db)
		page.Departements = model.AllDepartements(db)
		page.Recolteurs = model.AllRecolteurs(db)
	})

	p.On("form", "form", func(s *ihui.Session, event ihui.Event) {
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
	})

	p.On("input", "[id=search]", func(s *ihui.Session, event ihui.Event) {
		val := event.Value()
		log.Println(val)
		page.Individu.Latitude, page.Individu.Longitude, _ = model.FindLatLng(val)
		s.Script("updateEditMap(%f,%f)", page.Individu.Latitude, page.Individu.Longitude)
	})

	p.On("click", "[id=cancel]", func(s *ihui.Session, event ihui.Event) {
		s.QuitPage()
	})

	p.On("click", "[id=edit]", func(s *ihui.Session, event ihui.Event) {
		page.Edit = true
	})

	p.On("click", "[id=confirm-delete]", func(s *ihui.Session, event ihui.Event) {
		if page.Individu.ID > 0 {
			if err := db.Delete(page.Individu).Error; err != nil {
				log.Println(err)
				page.Error = err.Error()
				return
			}
		}
		s.QuitPage()
	})

	p.On("click", "[id=cancel-delete]", func(s *ihui.Session, event ihui.Event) {
		page.Delete = false
	})

	p.On("click", "[id=delete]", func(s *ihui.Session, event ihui.Event) {
		page.Delete = false
	})

	p.On("click", "[id=add-espece]", func(s *ihui.Session, event ihui.Event) {
		espece := model.Espece{}
		s.ShowPage(newPageEspece(&espece), &ihui.Options{Modal: true})
		if !db.NewRecord(espece) {
			page.Individu.Espece = espece
			page.Individu.EspeceID = espece.ID
		}
	})

	p.On("click", "[id=validation]", func(s *ihui.Session, event ihui.Event) {
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
		s.QuitPage()
	})

	p.On("position", "page", func(s *ihui.Session, event ihui.Event) {
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
	})

}
