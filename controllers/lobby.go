package controllers

import (
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
	this.Data["games"] = models.GetGameList([]int{models.StatusActive, models.StatusWait}, user.Id)

	this.LayoutSections = make(map[string]string)
	this.LayoutSections["Header"] = "components/navbar.html"
}
