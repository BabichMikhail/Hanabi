package controllers

import (
	"strconv"

	"github.com/BabichMikhail/Hanabi/models"
	"github.com/beego/wetalk/modules/auth"
)

type GameController struct {
	BaseController
}

func (this *GameController) Game() {
	id, _ := strconv.Atoi(this.Ctx.Input.Param(":id"))
	game, err := models.ReadGameById(id)
	if err != nil {
		gamePlayers := models.GetGamePlayers([]int{id})[id]
		playerIds := []int{}
		for i := 0; i < len(gamePlayers); i++ {
			playerIds = append(playerIds, gamePlayers[i].Id)
		}
		game, _ = models.CreateActiveGame(playerIds, id)
	}

	if game.IsGameOver() {
		this.Ctx.Redirect(302, this.URLFor("GameController.GameInactive", ":id", id))
	}

	userId := auth.GetUserIdFromSession(this.Ctx.Input.CruSession)
	playerInfo := game.GetPlayerGameInfo(userId)
	this.Data["playerInfo"] = playerInfo
	this.Data["S"] = game.SprintGame()
	this.Layout = "base.tpl"
	this.TplName = "templates/game.html"
}

func (this *GameController) GameInactive() {
	id, _ := strconv.Atoi(this.Ctx.Input.Param(":id"))

	game, _ := models.ReadGameById(id)
	if !game.IsGameOver() {
		this.Ctx.Redirect(302, this.URLFor("GameController.Game", ":id", id))
	}
	this.Data["Points"], _ = game.GetPoints()
	this.Layout = "base.tpl"
	this.TplName = "templates/gameinactive.html"
}
