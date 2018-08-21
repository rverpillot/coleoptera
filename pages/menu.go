package pages

import (
	"io"
	"rverpi/ihui"
)

type Menu struct {
	ihui.Component
	MenuItemActive string
	Connected      bool
	OnMenu         ihui.ClickAction
	OnConnection   ihui.ClickAction
	OnDisconnect   ihui.ClickAction
	PageIndividus  *PageIndividus `inject:""`
	PageEspeces    *PageEspeces   `inject:""`
	PagePlan       *PagePlan      `inject:""`
}

func NewMenu(menuItemActive string) *Menu {
	return &Menu{MenuItemActive: menuItemActive}
}

func (menu *Menu) OnInit(ctx *ihui.Context) {
	menu.PageEspeces.On("search_espece", func(ctx *ihui.Context) {
		ctx.Set("search_espece", ctx.Event.Data)
		menu.OnMenu("individus")
		menu.PagePlan.Trigger("refresh", ctx)
	})

	menu.PageIndividus.On("searching", func(ctx *ihui.Context) {
		menu.PagePlan.Trigger("refresh", ctx)
	})

	menu.OnMenu = func(menuItemActive string) {
		menu.MenuItemActive = menuItemActive
		switch menu.MenuItemActive {
		case "individus":
			ctx.DisplayPage(menu.PageIndividus, false)
		case "especes":
			ctx.DisplayPage(menu.PageEspeces, false)
		case "plan":
			ctx.DisplayPage(menu.PagePlan, false)
		}
	}

	menu.OnConnection = func(_ string) {
		ctx.DisplayPage(NewPageLogin(), true)
	}

	menu.OnDisconnect = func(_ string) {
		ctx.Set("admin", false)
	}
}

func (menu *Menu) Render(w io.Writer, ctx *ihui.Context) {
	menu.Connected = ctx.Get("admin").(bool)

	renderTemplate("menu", w, menu)
}
