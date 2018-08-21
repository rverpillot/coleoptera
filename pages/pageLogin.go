package pages

import (
	"io"
	"rverpi/ihui"
)

type PageLogin struct {
	*Page
	Error        string
	SubmitAction ihui.SubmitAction
	CancelAction ihui.ClickAction
}

func NewPageLogin() *PageLogin {
	return &PageLogin{Page: NewPage("login", true)}
}

func (page *PageLogin) OnInit(ctx *ihui.Context) {
	page.SubmitAction = func(data map[string]interface{}) {
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
			ctx.Set("admin", true)
			page.Close()
		} else {
			page.Error = "Utilisateur ou mot de passe inconnu!"
		}
	}
	page.CancelAction = func(_ string) {
		page.Close()
	}
}

func (page *PageLogin) Render(w io.Writer, ctx *ihui.Context) {
	page.renderPage(w, page)
}

func (page *PageLogin) authenticate(username string, password string) bool {
	if username == "admin" && password == "longicornes" {
		return true
	}
	return false
}
