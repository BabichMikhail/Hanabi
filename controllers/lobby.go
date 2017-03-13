package controllers

import (
	"github.com/BabichMikhail/Hanabi/models"
	"github.com/beego/wetalk/modules/auth"
	wetalk "github.com/beego/wetalk/modules/models"
)

type LobbyController struct {
	BaseController
}

func (c *LobbyController) GameList() {
	c.SetBaseLayout()
	c.TplName = "templates/gamelist.html"
	var user wetalk.User
	auth.GetUserFromSession(&user, c.Ctx.Input.CruSession)
	c.Data["user"] = user
	c.Data["games"] = models.GetGameList([]int{models.StatusActive, models.StatusWait}, user.Id)

	c.LayoutSections = make(map[string]string)
	c.LayoutSections["Header"] = "components/navbar.html"
	c.LayoutSections["Scripts"] = "scripts/lobbyscripts.tpl"
}
