package pages

import (
	"encoding/json"
	"io"
	"rverpi/coleoptera/model"
	"rverpi/ihui"

	"github.com/jinzhu/gorm"
)

type infoMap struct {
	Lat  float64 `json:"lat"`
	Lng  float64 `json:"lng"`
	Zoom int
}

func showMarkers(ctx *ihui.Context) {
	db := ctx.Get("db").(*gorm.DB)
	var espece_id uint
	var search string
	if ctx.Get("search_individus") != nil {
		search = ctx.Get("search_individus").(string)
	}
	if ctx.Get("search_espece") != nil {
		espece_id = ctx.Get("search_espece").(uint)
	}
	var markers []model.Marker
	markers, _ = model.Markers(db, search, espece_id)
	js, _ := json.Marshal(&markers)
	ctx.Script("showMarkers('#map',%s)", string(js))
}

type PagePlan struct {
	ihui.Page
	menu    *Menu
	infoMap infoMap
	refresh bool
}

func NewPagePlan(menu *Menu) *PagePlan {
	page := &PagePlan{menu: menu}
	page.Add("#menu", menu)
	return page
}

func (page *PagePlan) OnInit(ctx *ihui.Context) {
	page.infoMap.Lat = 46.435317
	page.infoMap.Lng = 1.812990
	page.infoMap.Zoom = 6

	page.On("map-loaded", func(ctx *ihui.Context) {
		showMarkers(ctx)

		data := ctx.Event.Data.(map[string]interface{})
		page.infoMap.Zoom = int(data["zoom"].(float64))
		page.infoMap.Lat = data["lat"].(float64)
		page.infoMap.Lng = data["lng"].(float64)
	})

	page.On("refresh", func(ctx *ihui.Context) {
		page.refresh = true
	})
}

func (page *PagePlan) OnShow(ctx *ihui.Context) {
	if page.refresh {
		showMarkers(ctx)
	}
}

func (page *PagePlan) Render(w io.Writer, ctx *ihui.Context) {
	renderTemplate("plan", w, page.infoMap)
}

func (page *PagePlan) OnDisplay(ctx *ihui.Context) {
	if page.refresh {
		// log.Printf("refreshMap({lat:%f, lng: %f}, %d)\n", page.infoMap.Lat, page.infoMap.Lng, page.infoMap.Zoom)
		ctx.Script("refreshMap({lat:%f, lng: %f}, %d)", page.infoMap.Lat, page.infoMap.Lng, page.infoMap.Zoom)
		page.refresh = false
	}
}
