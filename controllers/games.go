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
	gamePlayers := models.GetGamePlayers([]int{id})[id]
	playerIds := []int{}
	for i := 0; i < len(gamePlayers); i++ {
		playerIds = append(playerIds, gamePlayers[i].Id)
	}
	game, err := models.ReadActiveGameById(id)
	if err != nil {
		game, _ = models.CreateActiveGame(playerIds, id)
	}
	userId := auth.GetUserIdFromSession(this.Ctx.Input.CruSession)
	playerInfo := game.GetPlayerGameInfo(userId)
	this.Data["playerInfo"] = playerInfo
	this.Data["S"] = game.SprintGame()
	this.Layout = "base.tpl"
	this.TplName = "templates/game.html"
}
