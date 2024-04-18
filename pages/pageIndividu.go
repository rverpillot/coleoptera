package pages

import (
	"database/sql"
	"log"
	"strconv"
	"time"

	"github.com/jinzhu/gorm"
	"github.com/rverpillot/coleoptera/model"
	"github.com/rverpillot/ihui"
	"github.com/rverpillot/ihui/templating"
)

type PageIndividu struct {
	tmpl           *templating.PageAce
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

func (page *PageIndividu) Render(p *ihui.Page) error {
	db := p.Get("db").(*gorm.DB)
	page.Especes = model.AllEspeces(db)

	page.Admin = p.Get("admin").(bool)
	page.Sites = model.AllSites(db)
	page.Communes = model.AllCommunes(db)
	page.Departements = model.AllDepartements(db)
	page.Recolteurs = model.AllRecolteurs(db)

	if err := page.tmpl.Render(p); err != nil {
		return err
	}

	p.On("created", "page", func(s *ihui.Session, event ihui.Event) error {
		if page.Edit {
			return s.Script(`createEditMap("#mapedit","%s")`, p.Id)
		} else {
			return s.Script(`createPreviewMap("#mappreview",%f,%f)`, page.Individu.Longitude, page.Individu.Latitude)
		}
	})

	p.On("updated", "page", func(s *ihui.Session, event ihui.Event) error {
		if page.Edit && !page.EditMapCreated {
			if err := s.Script(`createEditMap("#mapedit","%s")`, p.Id); err != nil {
				return err
			}
		}
		page.EditMapCreated = true
		return nil
	})

	p.On("form", "form", func(s *ihui.Session, event ihui.Event) error {
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
			// lat, lng, err := model.FindLatLng(page.Individu.Commune)
			// if err != nil {
			// 	log.Println(err)
			// 	break
			// }
			// page.Individu.Latitude = lat
			// page.Individu.Longitude = lng
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
			page.Individu.Commentaire = sql.NullString{String: val, Valid: true}
		}
		return nil
	})

	p.On("click", "#cancel", func(s *ihui.Session, event ihui.Event) error {
		return p.Close()
	})

	p.On("click", "#edit", func(s *ihui.Session, event ihui.Event) error {
		page.Edit = true
		return nil
	})

	p.On("click", "#delete", func(s *ihui.Session, event ihui.Event) error {
		page.Delete = true
		return nil
	})

	p.On("click", "#confirm-delete", func(s *ihui.Session, event ihui.Event) error {
		if page.Individu.ID > 0 {
			if err := db.Delete(page.Individu).Error; err != nil {
				page.Error = err.Error()
				log.Println(err)
			}
		}
		return p.Close()
	})

	p.On("click", "#cancel-delete", func(s *ihui.Session, event ihui.Event) error {
		return p.Close()
	})

	p.On("click", "#add-espece", func(s *ihui.Session, event ihui.Event) error {
		espece := model.Espece{}
		s.ShowPage("espece", newPageEspece(&espece), &ihui.Options{Modal: true})
		if !db.NewRecord(espece) {
			page.Individu.Espece = espece
			page.Individu.EspeceID = espece.ID
		}
		return nil
	})

	p.On("click", "#validation", func(s *ihui.Session, event ihui.Event) error {
		log.Println(page.Individu)
		if page.Individu.Espece.ID == 0 {
			page.Error = "Genre/espèce absent !"
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
		}
		return p.Close()
	})

	p.On("position", "page", func(s *ihui.Session, event ihui.Event) error {
		pos := event.Data.(map[string]interface{})
		log.Println(pos)
		page.Individu.Longitude = pos["lng"].(float64)
		page.Individu.Latitude = pos["lat"].(float64)

		var altitude int64
		var err error
		page.Individu.Commune, page.Individu.Code, altitude, err = model.FindLocation(page.Individu.Longitude, page.Individu.Latitude)
		if err != nil {
			log.Println(err)
		}
		page.Individu.Altitude = sql.NullInt64{Int64: altitude, Valid: true}
		return nil
	})

	return nil
}
