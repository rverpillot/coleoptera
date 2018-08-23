package pages

import (
	"rverpi/ihui.v2"
)

type PageLogin struct {
	tmpl  *ihui.PageAce
	Error string
}

func NewPageLogin() *PageLogin {
	page := &PageLogin{}
	page.tmpl = newAceTemplate("login.ace", page)
	return page
}

func (page *PageLogin) Render(p ihui.Page) {
	page.tmpl.Render(p)

	p.On("submit", "form", func(s *ihui.Session, event ihui.Event) {
		data := event.Data.(map[string]interface{})
		username := data["username"].(string)
		password := data["password"].(string)
		if username == "" {
			page.Error = "Le nom d'utilisateur est vide!"
			return
		}
		if password == "" {
			page.Error = "Le mot de passe est vide!"
			return
		}
		if page.authenticate(username, password) {
			s.Set("admin", true)
			s.QuitPage()
		} else {
			page.Error = "Utilisateur ou mot de passe inconnu!"
		}
	})

	p.On("click", "[id=cancel]", func(s *ihui.Session, event ihui.Event) {
		s.QuitPage()
	})
}

func (page *PageLogin) authenticate(username string, password string) bool {
	if username == "admin" && password == "longicornes" {
		return true
	}
	return false
}
