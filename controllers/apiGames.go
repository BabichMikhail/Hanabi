package controllers

import (
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
		Status string                           `json:"status"`
		Colors map[gamePackage.CardColor]string `json:"colors"`
		Values map[gamePackage.CardValue]string `json:"values"`
	}{StatusSuccess, card.GetColors(), card.GetValues()}
	this.SetData(&result)
}

func (this *ApiGameController) GamePlayCard() {
	gameId, _ := this.GetInt("game_id")
	game, err := models.ReadActiveGameById(gameId)
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
	game, err := models.ReadActiveGameById(gameId)
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
	game, err := models.ReadActiveGameById(gameId)
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
	game, err := models.ReadActiveGameById(gameId)
	if this.SetError(err) {
		return
	}

	playerPosition, _ := this.GetInt("player_position")
	cardColor, _ := this.GetInt("card_color")
	game.NewActionInformationColor(playerPosition, gamePackage.CardColor(cardColor))
	models.UpdateCurrentGameById(gameId, game)
	this.SetSuccessResponse()
}

func (this *ApiGameController) GameCurrentStep() {
	gameId, _ := this.GetInt("game_id")
	game, err := models.ReadActiveGameById(gameId)
	if this.SetError(err) {
		return
	}
	result := struct {
		Status string `json:"status"`
		Step   int    `json:"step"`
	}{StatusSuccess, len(game.Actions)}
	this.SetData(&result)
}

func (this *ApiGameController) GameInfo() {
	gameId, _ := this.GetInt("game_id")
	game, err := models.ReadActiveGameById(gameId)
	if this.SetError(err) {
		return
	}

	playerPosition, err := game.GetPlayerPositionById(auth.GetUserIdFromSession(this.Ctx.Input.CruSession))
	if this.SetError(err) {
		return
	}

	result := struct {
		Status         string `json:"status"`
		Count          int    `json:"player_count"`
		PlayerPosition int    `json:"player_position"`
	}{StatusSuccess, game.PlayerCount, playerPosition}
	this.SetData(&result)
}
