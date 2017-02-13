package controllers

import (
	"strconv"

	"github.com/BabichMikhail/Hanabi/models"
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

	id, _ := strconv.Atoi(this.Ctx.Input.Param(":id"))
	game, err := models.ReadInactiveGameById(id)
	if err != nil {
		this.Ctx.Redirect(302, this.URLFor("LobbyController.GameList"))
	}
	this.Data["Game"] = game
	this.LayoutSections = make(map[string]string)
	this.LayoutSections["Header"] = "components/navbar.html"
	this.LayoutSections["Scripts"] = "components/viewscripts.html"
}
