package controllers

import (
	"strconv"

	gamePackage "github.com/BabichMikhail/Hanabi/engine/game"
	"github.com/BabichMikhail/Hanabi/models"
	"github.com/beego/wetalk/modules/auth"
)

type ApiGameController struct {
	ApiController
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
	this.SetData(&result)
}

func (this *ApiGameController) GetGameStatuses() {
	id, _ := strconv.Atoi(this.Ctx.Input.Param(":id"))
	gamePlayers := models.GetGamePlayers([]int{id})[id]
	playerIds := []int{}
	for i := 0; i < len(gamePlayers); i++ {
		playerIds = append(playerIds, gamePlayers[i].Id)
	}

	game, err := models.ReadGameById(id)
	if err != nil {
		game, _ = models.CreateActiveGame(playerIds, id)
	}
	userId := auth.GetUserIdFromSession(this.Ctx.Input.CruSession)
	playerInfo := game.GetPlayerGameInfo(userId)
	this.SetData(&playerInfo)
}

func (this *ApiGameController) GamePlayCard() {
	gameId, _ := this.GetInt("game_id")
	game, err := models.ReadGameById(gameId)
	if this.SetError(err) {
		return
	}

	playerPosition, err := game.GetPlayerPositionById(auth.GetUserIdFromSession(this.Ctx.Input.CruSession))
	if this.SetError(err) {
		return
	}

	cardPosition, _ := this.GetInt("card_position")
	game.NewActionPlaying(playerPosition, cardPosition)
	models.UpdateCurrentGameById(gameId, game)
	this.SetSuccessResponse()
}

func (this *ApiGameController) GameDiscardCard() {
	gameId, _ := this.GetInt("game_id")
	game, err := models.ReadGameById(gameId)
	if this.SetError(err) {
		return
	}

	playerPosition, err := game.GetPlayerPositionById(auth.GetUserIdFromSession(this.Ctx.Input.CruSession))
	if this.SetError(err) {
		return
	}

	cardPosition, _ := this.GetInt("card_position")
	game.NewActionDiscard(playerPosition, cardPosition)
	models.UpdateCurrentGameById(gameId, game)
	this.SetSuccessResponse()
}

func (this *ApiGameController) GameInfoCardValue() {
	gameId, _ := this.GetInt("game_id")
	game, err := models.ReadGameById(gameId)
	if this.SetError(err) {
		return
	}

	playerPosition, _ := this.GetInt("player_position")
	cardValue, _ := this.GetInt("card_value")
	game.NewActionInformationValue(playerPosition, gamePackage.CardValue(cardValue))
	models.UpdateCurrentGameById(gameId, game)
	this.SetSuccessResponse()
}

func (this *ApiGameController) GameInfoCardColor() {
	gameId, _ := this.GetInt("game_id")
	game, err := models.ReadGameById(gameId)
	if this.SetError(err) {
		return
	}

	playerPosition, _ := this.GetInt("player_position")
	cardColor, _ := this.GetInt("card_color")
	game.NewActionInformationColor(playerPosition, gamePackage.CardColor(cardColor))
	models.UpdateCurrentGameById(gameId, game)
	this.SetSuccessResponse()
}
