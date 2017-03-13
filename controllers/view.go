package controllers

import (
	"strconv"

	engineGame "github.com/BabichMikhail/Hanabi/game"
	"github.com/BabichMikhail/Hanabi/models"
	"github.com/beego/wetalk/modules/auth"
	wetalk "github.com/beego/wetalk/modules/models"
)

type GameViewController struct {
	BaseController
}

type CardUrl struct {
	Color engineGame.CardColor
	Value engineGame.CardValue
	Url   string
}

func (c *GameViewController) GameView() {
	c.SetBaseLayout()
	c.TplName = "templates/gameview.html"
	var user wetalk.User
	auth.GetUserFromSession(&user, c.Ctx.Input.CruSession)
	c.Data["user"] = user

	id, _ := strconv.Atoi(c.Ctx.Input.Param(":id"))
	var err error
	c.Data["InitState"], err = models.ReadInitialGameState(id)
	if err != nil {
		c.Ctx.Redirect(302, c.URLFor("LobbyController.GameList"))
	}

	c.Data["Actions"], err = models.ReadActions(id)
	if err != nil {
		c.Ctx.Redirect(302, c.URLFor("LobbyController.GameList"))
	}

	c.LayoutSections = make(map[string]string)
	c.LayoutSections["Header"] = "components/navbar.html"
	c.LayoutSections["Scripts"] = "scripts/viewscripts.tpl"

	var urls []CardUrl
	for _, color := range engineGame.Colors {
		for _, value := range engineGame.Values {
			urls = append(urls, CardUrl{
				Color: color,
				Value: value,
				Url:   engineGame.GetCardUrlByValueAndColor(color, value),
			})
		}
	}

	c.Data["Players"] = models.GetGamePlayers([]int{id})[id]
	c.Data["CardUrls"] = urls
	c.Data["MaxRedTokens"] = engineGame.MaxRedTokens
	c.Data["MaxBlueTokens"] = engineGame.MaxBlueTokens
	c.Data["NoneColor"] = engineGame.NoneColor
	c.Data["NoneValue"] = engineGame.NoneValue
	c.Data["TableColors"] = engineGame.GetTableColorOrder()
	c.Data["ActionTypes"] = map[string]int{
		"infoColor": engineGame.TypeActionInformationColor,
		"infoValue": engineGame.TypeActionInformationValue,
		"discard":   engineGame.TypeActionDiscard,
		"play":      engineGame.TypeActionPlaying,
	}
}
