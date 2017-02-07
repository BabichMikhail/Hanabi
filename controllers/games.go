package controllers

import (
	"strconv"

	engineGame "github.com/BabichMikhail/Hanabi/engine/game"
	"github.com/BabichMikhail/Hanabi/models"
	"github.com/beego/wetalk/modules/auth"
	wetalk "github.com/beego/wetalk/modules/models"
)

type GameController struct {
	BaseController
}

func (this *GameController) Game() {
	id, _ := strconv.Atoi(this.Ctx.Input.Param(":id"))
	game, err := models.ReadActiveGameById(id)

	if err != nil {
		this.Ctx.Redirect(302, this.URLFor("LobbyController.GameList"))
		return
	}

	if game.IsGameOver() {
		this.Ctx.Redirect(302, this.URLFor("GameController.GameInactive", ":id", id))
		return
	}

	userId := auth.GetUserIdFromSession(this.Ctx.Input.CruSession)
	playerInfo := game.GetPlayerGameInfo(userId)
	this.Data["playerInfo"] = playerInfo
	this.Data["deckFirstNumber"] = playerInfo.DeckSize / 10
	this.Data["deckSecondNumber"] = playerInfo.DeckSize % 10

	this.Data["Step"] = len(game.Actions)
	this.Layout = "base.tpl"
	this.TplName = "templates/game.html"

	var user wetalk.User
	auth.GetUserFromSession(&user, this.Ctx.Input.CruSession)
	this.Data["user"] = user
	this.LayoutSections = make(map[string]string)
	this.LayoutSections["Header"] = "components/navbar.html"
	this.Data["TableColors"] = engineGame.GetTableColorOrder()
	playerNickNames := []string{}
	for i := 0; i < len(game.CurrentState.PlayerStates); i++ {
		nickName := models.GetUserNickNameById(game.CurrentState.PlayerStates[i].PlayerId)
		playerNickNames = append(playerNickNames, nickName)
	}
	this.Data["NickNames"] = playerNickNames
}

func (this *GameController) GameInactive() {
	id, _ := strconv.Atoi(this.Ctx.Input.Param(":id"))

	game, err := models.ReadGameById(id)
	if err != nil {
		this.Ctx.Redirect(302, this.URLFor("LobbyController.GameList"))
		return
	}
	if !game.IsGameOver() {
		this.Ctx.Redirect(302, this.URLFor("GameController.Game", ":id", id))
		return
	}
	this.Data["Points"], _ = game.GetPoints()
	this.Layout = "base.tpl"
	this.TplName = "templates/gameinactive.html"

	var user wetalk.User
	auth.GetUserFromSession(&user, this.Ctx.Input.CruSession)
	this.Data["user"] = user
	this.LayoutSections = make(map[string]string)
	this.LayoutSections["Header"] = "components/navbar.html"
}
