package pages

import (
	"fmt"

	"github.com/rverpillot/ihui"
)

type Item struct {
	Name   string
	Label  string
	Active bool
	Drawer ihui.HTMLRenderer
}

type Menu struct {
	Connected bool
	Items     []*Item
}

func NewMenu() *Menu {
	return &Menu{}
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

func (menu *Menu) ShowItem(s *ihui.Session, name string) error {
	for _, item := range menu.Items {
		if item.Name == name {
			menu.SetActive(name)
			if err := s.ShowPage(item.Name, item.Drawer); err != nil {
				fmt.Println(err)
				return err
			}
			break
		}
	}
	return nil
}

func (menu *Menu) Render(e *ihui.HTMLElement) error {
	menu.Connected = e.Get("admin").(bool)

	e.OnClick(".menu-item", func(s *ihui.Session, event ihui.Event) error {
		return menu.ShowItem(s, event.Value())
	})

	e.OnClick("#connect", func(s *ihui.Session, _ ihui.Event) error {
		return s.ShowModal("login", NewPageLogin())
	})

	e.OnClick("#disconnect", func(s *ihui.Session, _ ihui.Event) error {
		s.Set("admin", false)
		return nil
	})

	return e.WriteGoTemplate(TemplatesFs, "templates/menu.html", menu)
}
