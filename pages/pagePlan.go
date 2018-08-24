package pages

import (
	"encoding/json"

	"rverpi/coleoptera.v3/model"
	"rverpi/ihui.v2"

	"github.com/jinzhu/gorm"
)

type infoMap struct {
	Lat  float64 `json:"lat"`
	Lng  float64 `json:"lng"`
	Zoom int
}

type PagePlan struct {
	tmpl    *ihui.PageAce
	menu    *Menu
	infoMap infoMap
	refresh bool
}

func NewPagePlan() *PagePlan {
	page := &PagePlan{
		infoMap: infoMap{
			Lat:  46.435317,
			Lng:  1.812990,
			Zoom: 6,
		},
	}
	page.tmpl = newAceTemplate("plan.ace", page)
	return page
}

func (page *PagePlan) Render(p ihui.Page) {
	page.tmpl.Render(p)

	if page.refresh {
		p.On("updated", "page", func(s *ihui.Session, event ihui.Event) {
			page.showMarkers(p.Session())
			p.Script("refreshMap({lat:%f, lng: %f}, %d)", page.infoMap.Lat, page.infoMap.Lng, page.infoMap.Zoom)
		})
		page.refresh = false
	}

	p.On("load", "page", func(s *ihui.Session, event ihui.Event) {
		s.Script(`createMap("#map", {lat:%f, lng:%f}, %d)`, page.infoMap.Lat, page.infoMap.Lng, page.infoMap.Zoom)
		page.showMarkers(s)
	})

	p.On("map-changed", "page", func(s *ihui.Session, event ihui.Event) {
		data := event.Data.(map[string]interface{})
		page.infoMap = infoMap{
			Lat:  data["lat"].(float64),
			Lng:  data["lng"].(float64),
			Zoom: int(data["zoom"].(float64)),
		}
	})

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

func (page *PagePlan) RefreshMarkers() {
	page.refresh = true
}
