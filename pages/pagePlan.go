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
	page := &PagePlan{}
	page.tmpl = newAceTemplate("plan.ace", page)
	return page
}

func (page *PagePlan) Render(p ihui.Page) {
	page.tmpl.Render(p)

	if page.refresh {
		// p.Script("refreshMap({lat:%f, lng: %f}, %d)", page.infoMap.Lat, page.infoMap.Lng, page.infoMap.Zoom)
		page.showMarkers(p.Session())
		page.refresh = false
	}

	p.On("geoportail-ready", "page", func(s *ihui.Session, event ihui.Event) {
		page.infoMap.Lat = 46.435317
		page.infoMap.Lng = 1.812990
		page.infoMap.Zoom = 6
		s.Script(`createMap("#map", {lat:%f, lng:%f}, %d)`, page.infoMap.Lat, page.infoMap.Lng, page.infoMap.Zoom)
		page.showMarkers(s)
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
