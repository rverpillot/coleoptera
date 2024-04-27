package pages

import (
	"encoding/json"

	"github.com/rverpillot/coleoptera/model"
	"github.com/rverpillot/ihui"
	"gorm.io/gorm"
)

type infoMap struct {
	Lat  float64 `json:"lat"`
	Lng  float64 `json:"lng"`
	Zoom int
}

type PagePlan struct {
	menu    *Menu
	infoMap infoMap
}

func NewPagePlan(menu *Menu) *PagePlan {
	return &PagePlan{
		menu: menu,
		infoMap: infoMap{
			Lat:  0,
			Lng:  0,
			Zoom: 8,
		},
	}
}

func (page *PagePlan) Render(e *ihui.HTMLElement) error {
	if err := e.WriteGoTemplate(TemplatesFs, "templates/plan.html", page); err != nil {
		return err
	}

	e.On("element-created", "", func(s *ihui.Session, event ihui.Event) error {
		s.Execute(`createMap("#map", {lat:%f, lng:%f}, %d, "%s")`,
			page.infoMap.Lat, page.infoMap.Lng, page.infoMap.Zoom, e.Id)
		page.showMarkers(s)
		return nil
	})

	e.On("element-updated", "", func(s *ihui.Session, event ihui.Event) error {
		page.showMarkers(s)
		return nil
	})

	e.On("map-changed", "", func(s *ihui.Session, event ihui.Event) error {
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
	session.Execute("showMarkers('#map',%s)", js)
}
