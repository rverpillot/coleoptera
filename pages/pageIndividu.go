package pages

import (
	"database/sql"
	"log"
	"strconv"
	"time"

	"github.com/rverpillot/coleoptera/model"
	"github.com/rverpillot/ihui"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type PageIndividu struct {
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
	return &PageIndividu{
		Individu: individu,
		Edit:     editMode,
	}
}

func (page *PageIndividu) Render(e *ihui.HTMLElement) error {
	db := e.Get("db").(*gorm.DB)
	page.Especes = model.AllEspeces(db)

	page.Admin = e.Get("admin").(bool)
	page.Sites = model.AllSites(db)
	page.Communes = model.AllCommunes(db)
	page.Departements = model.AllDepartements(db)
	page.Recolteurs = model.AllRecolteurs(db)

	e.On("element-created", "", func(s *ihui.Session, event ihui.Event) error {
		if page.Edit {
			return s.Execute(`createEditMap("#mapedit","%s")`, e.Id)
		} else {
			return s.Execute(`createPreviewMap("#mappreview",%f,%f)`, page.Individu.Longitude, page.Individu.Latitude)
		}
	})

	e.On("element-updated", "", func(s *ihui.Session, event ihui.Event) error {
		if page.Edit && !page.EditMapCreated {
			if err := s.Execute(`createEditMap("#mapedit","%s")`, e.Id); err != nil {
				return err
			}
		}
		page.EditMapCreated = true
		return nil
	})

	e.OnForm("form", func(s *ihui.Session, event ihui.Event) error {
		data := event.Data.(map[string]interface{})
		name := data["name"].(string)
		val := data["val"].(string)
		log.Printf("%s=%s", name, val)

		switch name {
		case "date":
			page.Individu.Date, _ = time.Parse("2006-01-02", val)
		case "espece":
			id, _ := strconv.Atoi(val)
			page.Individu.EspeceID = uint(id)
			page.Individu.Espece.ID = uint(id)
			db.First(&page.Individu.Espece)
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
			s.Execute("updateEditMap(%f,%f)", page.Individu.Latitude, page.Individu.Longitude)
		case "latitude":
			lat, _ := strconv.ParseFloat(val, 64)
			page.Individu.Latitude = lat
			page.Search = ""
			s.Execute("updateEditMap(%f,%f)", page.Individu.Latitude, page.Individu.Longitude)
		case "recolteur":
			page.Individu.Recolteur = val
		case "commentaire":
			page.Individu.Commentaire = sql.NullString{String: val, Valid: true}
		}
		return nil
	})

	e.OnClick("#cancel", func(s *ihui.Session, event ihui.Event) error {
		return e.Close()
	})

	e.OnClick("#delete", func(s *ihui.Session, event ihui.Event) error {
		page.Delete = true
		return nil
	})

	e.OnClick("#confirm-delete", func(s *ihui.Session, event ihui.Event) error {
		if page.Individu.ID > 0 {
			if err := db.Delete(page.Individu).Error; err != nil {
				page.Error = err.Error()
				log.Println(err)
			}
		}
		return e.Close()
	})

	e.OnClick("#cancel-delete", func(s *ihui.Session, event ihui.Event) error {
		return e.Close()
	})

	e.OnClick("#validation", func(s *ihui.Session, event ihui.Event) error {
		log.Println(page.Individu)
		if page.Individu.Espece.ID == 0 {
			page.Error = "Genre/espèce absent !"
			return nil
		}
		if err := db.Omit(clause.Associations).Save(&page.Individu).Error; err != nil {
			log.Println(err)
			page.Error = err.Error()
			return nil
		}
		return e.Close()
	})

	e.On("position", "", func(s *ihui.Session, event ihui.Event) error {
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

	return e.WriteGoTemplate(TemplatesFs, "templates/individu.html", page)
}
