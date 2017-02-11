package controllers

import (
	"github.com/beego/wetalk/modules/auth"
	wetalk "github.com/beego/wetalk/modules/models"
)

type GameViewController struct {
	BaseController
}

func (this *GameViewController) GameView() {
	this.Layout = "base.tpl"
	this.TplName = "templates/gameview.html"
	var user wetalk.User
	auth.GetUserFromSession(&user, this.Ctx.Input.CruSession)
	this.Data["user"] = user

	this.LayoutSections = make(map[string]string)
	this.LayoutSections["Header"] = "components/navbar.html"
}
