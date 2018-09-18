package pages

import (
	"bitbucket.org/rverpi90/ihui"
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

	p.On("submit", "form", func(s *ihui.Session, event ihui.Event) bool {
		data := event.Data.(map[string]interface{})
		username := data["username"].(string)
		password := data["password"].(string)
		if username == "" {
			page.Error = "Le nom d'utilisateur est vide!"
			return true
		}
		if password == "" {
			page.Error = "Le mot de passe est vide!"
			return true
		}
		if page.authenticate(username, password) {
			s.Set("admin", true)
			s.QuitPage()
		} else {
			page.Error = "Utilisateur ou mot de passe inconnu!"
		}
		return true
	})

	p.On("click", "#cancel", func(s *ihui.Session, event ihui.Event) bool {
		return s.QuitPage()
	})
}

func (page *PageLogin) authenticate(username string, password string) bool {
	if username == "admin" && password == "longicornes" {
		return true
	}
	return false
}
