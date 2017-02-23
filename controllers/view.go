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

func (this *GameViewController) GameView() {
	this.SetBaseLayout()
	this.TplName = "templates/gameview.html"
	var user wetalk.User
	auth.GetUserFromSession(&user, this.Ctx.Input.CruSession)
	this.Data["user"] = user

	id, _ := strconv.Atoi(this.Ctx.Input.Param(":id"))
	var err error
	this.Data["InitState"], err = models.ReadInitialGameState(id)
	if err != nil {
		this.Ctx.Redirect(302, this.URLFor("LobbyController.GameList"))
	}

	this.Data["Actions"], err = models.ReadActions(id)
	if err != nil {
		this.Ctx.Redirect(302, this.URLFor("LobbyController.GameList"))
	}

	this.LayoutSections = make(map[string]string)
	this.LayoutSections["Header"] = "components/navbar.html"
	this.LayoutSections["Scripts"] = "components/viewscripts.html"

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

	this.Data["Players"] = models.GetGamePlayers([]int{id})[id]
	this.Data["CardUrls"] = urls
	this.Data["MaxRedTokens"] = engineGame.MaxRedTokens
	this.Data["MaxBlueTokens"] = engineGame.MaxBlueTokens
	this.Data["NoneColor"] = engineGame.NoneColor
	this.Data["NoneValue"] = engineGame.NoneValue
	this.Data["TableColors"] = engineGame.GetTableColorOrder()
	this.Data["ActionTypes"] = map[string]int{
		"infoColor": engineGame.TypeActionInformationColor,
		"infoValue": engineGame.TypeActionInformationValue,
		"discard":   engineGame.TypeActionDiscard,
		"play":      engineGame.TypeActionPlaying,
	}
}
