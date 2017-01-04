package controllers

import (
	"strconv"

	gamePackage "github.com/BabichMikhail/Hanabi/engine/game"
	"github.com/BabichMikhail/Hanabi/models"
	"github.com/beego/wetalk/modules/auth"
)

type ApiGameController struct {
	BaseController
}

var card gamePackage.Card

func init() {
	card = gamePackage.Card{}
}

func (this *ApiGameController) GetGameCards() {
	result := struct {
		Colors map[gamePackage.CardColor]string `json:"colors"`
		Values map[gamePackage.CardValue]string `json:"values"`
	}{card.GetColors(), card.GetValues()}
	this.Data["json"] = &result
	this.ServeJSON()
}

func (this *ApiGameController) GetGameStatuses() {
	id, _ := strconv.Atoi(this.Ctx.Input.Param(":id"))
	gamePlayers := models.GetGamePlayers([]int{id})[id]
	playerIds := []int{}
	for i := 0; i < len(gamePlayers); i++ {
		playerIds = append(playerIds, gamePlayers[i].Id)
	}
	game, err := models.ReadActiveGameByGameId(id)
	if err != nil {
		game, _ = models.CreateActiveGame(playerIds, id)
	}
	userId := auth.GetUserIdFromSession(this.Ctx.Input.CruSession)
	playerInfo := game.GetPlayerGameInfo(userId)
	this.Data["json"] = &playerInfo
	this.ServeJSON()
}
