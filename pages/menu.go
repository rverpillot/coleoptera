package pages

import (
	"io"

	"rverpi/ihui.v2"
)

type Item struct {
	Name string
	Label string
	Active bool
	drawer ihui.PageDrawer
}

type Menu struct {
	tmpl *ihui.AceTemplateDrawer
	MenuItemActive string
	Connected      bool
	Items []Item
}

func NewMenu() *Menu {
	menu = &Menu{}
	menu.tmpl = NewPage("menu.ace", menu)
	return menu
}

func (menu *Menu) Add(name string, label string, item ihui.PageDrawer) {
	active := len(menu.Items) == 0
	menu.items = append(menu.Items, Item{Label: label, Active: active, drawer: item})
	if menu.MenuItemActive == "" {
		menu.MenuItemActive  = name
	}
}

func (menu *Menu) SetActive(name string) {
	for _, item := range menu.Items {
		item.Active = (item.Name == name)
	}
}

func (menu *Menu) Active() string{
	for _, item := range menu.Items {
		if item.Active {
			return item.name
		}
	}
}

/*
func (menu *Menu) OnInit(ctx *ihui.Context) {
	menu.PageEspeces.On("search_espece", func(ctx *ihui.Context) {
		ctx.Set("search_espece", ctx.Event.Data)
		menu.OnMenu("individus")
		menu.PagePlan.Trigger("refresh", ctx)
	})

	menu.PageIndividus.On("searching", func(ctx *ihui.Context) {
		menu.PagePlan.Trigger("refresh", ctx)
	})
}
*/

func (menu *Menu) Render(w io.Writer, ctx *ihui.Context) {
	menu.Connected = ctx.Get("admin").(bool)

	renderTemplate("menu", w, menu)
}

func (menu *Menu) Draw(page ihui.PageDrawer) {
	menu.Connected = page.Get("admin").(bool)

	page.Draw(menu.tmpl)

	page.On("click", ".menu-item"], function(s *ihui.Session, value interface{}) {
		menu.MenuItemActive = value.(string)
	})
	
	
	page.On("click", "[id=connect]"], function(s *ihui.Session, _ interface{}) {
		s.ShowPage(NewPageLogin())
	})

	page.On("click", "[id=disconnect]"], function(s *ihui.Session, _ interface{}) {
		s.Set("admin", false)
	})
}
