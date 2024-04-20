package pages

import (
	"encoding/json"

	"github.com/jinzhu/gorm"
	"github.com/rverpillot/coleoptera/model"
	"github.com/rverpillot/ihui"
	"github.com/rverpillot/ihui/templating"
)

type infoMap struct {
	Lat  float64 `json:"lat"`
	Lng  float64 `json:"lng"`
	Zoom int
}

type PagePlan struct {
	tmpl    *templating.PageAce
	menu    *Menu
	infoMap infoMap
}

func NewPagePlan(menu *Menu) *PagePlan {
	return &PagePlan{
		tmpl: newAceTemplate("plan.ace"),
		menu: menu,
		infoMap: infoMap{
			Lat:  46.435317,
			Lng:  1.812990,
			Zoom: 5,
		},
	}
}

func (page *PagePlan) Render(p *ihui.Page) error {
	if err := p.WriteTemplate(page.tmpl, page); err != nil {
		return err
	}

	p.On("page-created", "", func(s *ihui.Session, event ihui.Event) error {
		s.Script(`createMap("#map", {lat:%f, lng:%f}, %d, "%s")`,
			page.infoMap.Lat, page.infoMap.Lng, page.infoMap.Zoom, p.Id)
		page.showMarkers(s)
		return nil
	})

	p.On("page-updated", "", func(s *ihui.Session, event ihui.Event) error {
		page.showMarkers(s)
		return nil
	})

	p.On("map-changed", "", func(s *ihui.Session, event ihui.Event) error {
		data := event.Data.(map[string]interface{})
		page.infoMap = infoMap{
			Lat:  data["lat"].(float64),
			Lng:  data["lng"].(float64),
			Zoom: int(data["zoom"].(float64)),
		}
		return nil
	})

	return nil
}

func (page *PagePlan) showMarkers(session *ihui.Session) {
	db := session.Get("db").(*gorm.DB)
	var espece_id uint
	var search string
	if session.Get("search_individus") != nil {
		search = session.Get("search_individus").(string)
	}
	if session.Get("search_espece") != nil {
		espece_id = session.Get("search_espece").(uint)
	}
	var markers []model.Marker
	markers, _ = model.Markers(db, search, espece_id)
	data, _ := json.Marshal(&markers)
	js := string(data)
	session.Script("showMarkers('#map',%s)", js)
}
