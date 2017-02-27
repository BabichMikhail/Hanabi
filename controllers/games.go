package controllers

import (
	"strconv"

	engineGame "github.com/BabichMikhail/Hanabi/game"
	"github.com/BabichMikhail/Hanabi/models"
	"github.com/beego/wetalk/modules/auth"
	wetalk "github.com/beego/wetalk/modules/models"
)

type GameController struct {
	BaseController
}

func (this *GameController) Game() {
	id, _ := strconv.Atoi(this.Ctx.Input.Param(":id"))
	state, err := models.ReadCurrentGameState(id)

	if err != nil {
		this.Ctx.Redirect(302, this.URLFor("LobbyController.GameList"))
		return
	}

	if state.IsGameOver() {
		models.SetGameFinishedStatus(id)
		this.Ctx.Redirect(302, this.URLFor("GameController.GameFinished", ":id", id))
		return
	}

	userId := auth.GetUserIdFromSession(this.Ctx.Input.CruSession)
	playerInfo := state.GetPlayerGameInfo(userId)
	this.Data["playerInfo"] = playerInfo
	this.Data["deckFirstNumber"] = playerInfo.DeckSize / 10
	this.Data["deckSecondNumber"] = playerInfo.DeckSize % 10

	this.Data["Step"], _ = models.GetActionCount(id)
	this.SetBaseLayout()
	this.TplName = "templates/game.html"

	var user wetalk.User
	auth.GetUserFromSession(&user, this.Ctx.Input.CruSession)
	this.Data["user"] = user
	this.LayoutSections = make(map[string]string)
	this.LayoutSections["Header"] = "components/navbar.html"
	this.Data["TableColors"] = engineGame.GetTableColorOrder()
	playerNickNames := []string{}
	for i := 0; i < len(state.PlayerStates); i++ {
		nickName := models.GetUserNickNameById(state.PlayerStates[i].PlayerId)
		playerNickNames = append(playerNickNames, nickName)
	}
	this.Data["NickNames"] = playerNickNames
}

func (this *GameController) GameFinished() {
	id, _ := strconv.Atoi(this.Ctx.Input.Param(":id"))

	state, err := models.ReadCurrentGameState(id)
	if err != nil {
		this.Ctx.Redirect(302, this.URLFor("LobbyController.GameList"))
		return
	}
	if !state.IsGameOver() {
		this.Ctx.Redirect(302, this.URLFor("GameController.Game", ":id", id))
		return
	}
	this.Data["Points"], _ = state.GetPoints()
	this.SetBaseLayout()
	this.TplName = "templates/gamefinished.html"

	var user wetalk.User
	auth.GetUserFromSession(&user, this.Ctx.Input.CruSession)
	this.Data["user"] = user
	this.LayoutSections = make(map[string]string)
	this.LayoutSections["Header"] = "components/navbar.html"
}
