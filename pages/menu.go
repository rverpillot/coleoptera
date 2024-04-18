package pages

import (
	"github.com/rverpillot/ihui"
	"github.com/rverpillot/ihui/templating"
)

type Item struct {
	Name   string
	Label  string
	Active bool
	Drawer ihui.HTMLRenderer
}

type Menu struct {
	tmpl      *templating.PageAce
	Connected bool
	Items     []*Item
}

func NewMenu() *Menu {
	menu := &Menu{}
	menu.tmpl = newAceTemplate("menu.ace", menu)
	return menu
}

func (menu *Menu) Add(name string, label string, drawer ihui.HTMLRenderer) {
	active := len(menu.Items) == 0
	menu.Items = append(menu.Items, &Item{Name: name, Label: label, Active: active, Drawer: drawer})
}

func (menu *Menu) SetActive(name string) {
	for _, item := range menu.Items {
		item.Active = (item.Name == name)
	}
}

func (menu *Menu) ShowPage(s *ihui.Session, name string) {
	for _, item := range menu.Items {
		if item.Name == name {
			menu.SetActive(name)
			s.ShowPage(item.Name, item.Drawer, &ihui.Options{Replace: true, Target: "#" + item.Name, Visible: true})
		} else {
			s.HidePage(item.Name)
		}
	}
}

func (menu *Menu) Render(page *ihui.Page) error {
	menu.Connected = page.Get("admin").(bool)

	if err := menu.tmpl.Render(page); err != nil {
		return err
	}

	page.On("click", ".menu-item", func(s *ihui.Session, event ihui.Event) {
		menu.ShowPage(s, event.Value())
	})

	page.On("click", "#connect", func(s *ihui.Session, _ ihui.Event) {
		s.ShowPage("login", NewPageLogin(), &ihui.Options{Modal: true})
	})

	page.On("click", "#disconnect", func(s *ihui.Session, _ ihui.Event) {
		s.Set("admin", false)
	})
	return nil
}
