package controllers

import (
	engine "github.com/BabichMikhail/Hanabi/engine"
	"github.com/BabichMikhail/Hanabi/models"
	"github.com/beego/wetalk/modules/auth"
	wetalk "github.com/beego/wetalk/modules/models"
)

type LobbyController struct {
	BaseController
}

func (this *LobbyController) GameList() {
	this.Layout = "base.tpl"
	this.TplName = "templates/gamelist.html"
	var user wetalk.User
	auth.GetUserFromSession(&user, this.Ctx.Input.CruSession)
	this.Data["user"] = user
	this.Data["games"] = models.GetGameList(engine.GetAllStatuses(), user.Id)
	if !this.Ctx.Input.IsPost() {
		return
	}
}
