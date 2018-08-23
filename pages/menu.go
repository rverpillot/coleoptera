package pages

import (
	"rverpi/ihui.v2"
)

type Item struct {
	Name   string
	Label  string
	Active bool
	drawer ihui.PageRenderer
}

type Menu struct {
	tmpl           *ihui.PageAce
	MenuItemActive string
	Connected      bool
	Items          []Item
}

func NewMenu() *Menu {
	menu := &Menu{}
	menu.tmpl = newAceTemplate("menu.ace", menu)
	return menu
}

func (menu *Menu) Add(name string, label string, item ihui.PageRenderer) {
	active := len(menu.Items) == 0
	menu.Items = append(menu.Items, Item{Label: label, Active: active, drawer: item})
	if menu.MenuItemActive == "" {
		menu.MenuItemActive = name
	}
}

func (menu *Menu) SetActive(name string) {
	for _, item := range menu.Items {
		item.Active = (item.Name == name)
	}
}

func (menu *Menu) Active() string {
	for _, item := range menu.Items {
		if item.Active {
			return item.Name
		}
	}
	return ""
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

func (menu *Menu) Render(page ihui.Page) {
	menu.Connected = page.Get("admin").(bool)

	menu.tmpl.Render(page)

	page.On("click", ".menu-item", func(s *ihui.Session, event ihui.Event) {
		menu.MenuItemActive = event.Value()
	})

	page.On("click", "#connect", func(s *ihui.Session, _ ihui.Event) {
		s.ShowPage(NewPageLogin(), nil)
	})

	page.On("click", "#disconnect", func(s *ihui.Session, _ ihui.Event) {
		s.Set("admin", false)
	})
}
