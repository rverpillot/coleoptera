package pages

import (
	"rverpi/ihui.v2"
)

type PageMain struct {
	tmpl          *ihui.PageAce
	Menu          *Menu
	pageEspeces   *PageEspeces
	pageIndividus *PageIndividus
	pagePlan      *PagePlan
}

func NewPageMain() *PageMain {
	page := &PageMain{
		Menu:          NewMenu(),
		pageEspeces:   NewPageEspeces(),
		pageIndividus: NewPageIndividus(),
		pagePlan:      NewPagePlan(),
	}
	page.tmpl = newAceTemplate("main.ace", page)
	page.Menu.Add("especes", "Espèces")
	page.Menu.Add("individus", "Individus")
	page.Menu.Add("plan", "Plan")
	page.Menu.SetActive("individus")
	return page
}

func (main *PageMain) Render(page ihui.Page) {
	main.tmpl.Render(page)
	page.Add("#menu", main.Menu)
	page.Add("#especes", main.pageEspeces)
	page.Add("#individus", main.pageIndividus)
	page.Add("#plan", main.pagePlan)

	page.On("click", ".espece", func(session *ihui.Session, _ ihui.Event) {
		main.Menu.SetActive("individus")
	})

	page.On("click", ".menu-item", func(s *ihui.Session, event ihui.Event) {
		if event.Value() == "plan" {
			main.pagePlan.RefreshMarkers()
		}
	})

}
