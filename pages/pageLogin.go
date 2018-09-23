package pages

import (
	"bytes"
	"crypto/sha256"

	"github.com/rverpillot/ihui"
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
			s.CloseModalPage()
		} else {
			page.Error = "Utilisateur ou mot de passe inconnu!"
		}
		return true
	})

	p.On("click", "#cancel", func(s *ihui.Session, event ihui.Event) bool {
		return s.CloseModalPage()
	})
}

func (page *PageLogin) authenticate(username string, password string) bool {
	ref := []byte{149, 247, 20, 30, 104, 99, 228, 222, 33, 243, 48, 132, 125, 204, 248, 211, 26, 247, 51, 254, 100, 182, 64, 47, 199, 119, 60, 197, 4, 127, 234, 167}
	h := sha256.New()
	h.Write([]byte(password))
	if username == "admin" && bytes.Equal(h.Sum(nil), ref) {
		return true
	}
	return false
}
