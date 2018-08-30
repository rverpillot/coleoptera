package pages

import (
	"rverpi90/ihui.v2"
)

type Item struct {
	Name   string
	Label  string
	Active bool
	Drawer ihui.PageRenderer
}

type Menu struct {
	tmpl      *ihui.PageAce
	Connected bool
	Items     []*Item
}

func NewMenu() *Menu {
	menu := &Menu{}
	menu.tmpl = newAceTemplate("menu.ace", menu)
	return menu
}

func (menu *Menu) Add(name string, label string, drawer ihui.PageRenderer) {
	active := len(menu.Items) == 0
	menu.Items = append(menu.Items, &Item{Name: name, Label: label, Active: active, Drawer: drawer})
}

func (menu *Menu) SetActive(name string) {
	for _, item := range menu.Items {
		item.Active = (item.Name == name)
	}
}

func (menu *Menu) ShowPage(s *ihui.Session, name string) bool {
	for _, item := range menu.Items {
		if item.Name == name {
			menu.SetActive(name)
			return s.ShowPage(item.Name, item.Drawer, nil)
		}
	}
	return false
}

func (menu *Menu) Render(page ihui.Page) {
	menu.Connected = page.Get("admin").(bool)

	menu.tmpl.Render(page)

	page.On("click", ".menu-item", func(s *ihui.Session, event ihui.Event) bool {
		return menu.ShowPage(s, event.Value())
	})

	page.On("click", "#connect", func(s *ihui.Session, _ ihui.Event) bool {
		return s.ShowPage("login", NewPageLogin(), &ihui.Options{Modal: true})
	})

	page.On("click", "#disconnect", func(s *ihui.Session, _ ihui.Event) bool {
		s.Set("admin", false)
		return true
	})

}
