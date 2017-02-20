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
	state, err := models.ReadCurrentGameState(gameId)
	if this.SetError(err) {
		return
	}

	playerPosition, err := state.GetPlayerPositionById(auth.GetUserIdFromSession(this.Ctx.Input.CruSession))
	if this.SetError(err) {
		return
	}

	cardPosition, _ := this.GetInt("card_position")
	action, _ := state.NewActionPlaying(playerPosition, cardPosition)
	models.NewAction(gameId, action)
	models.UpdateGameState(gameId, state)
	this.SetSuccessResponse()
}

func (this *ApiGameController) GameDiscardCard() {
	gameId, _ := this.GetInt("game_id")
	state, err := models.ReadCurrentGameState(gameId)
	if this.SetError(err) {
		return
	}

	playerPosition, err := state.GetPlayerPositionById(auth.GetUserIdFromSession(this.Ctx.Input.CruSession))
	if this.SetError(err) {
		return
	}

	cardPosition, _ := this.GetInt("card_position")
	action, _ := state.NewActionDiscard(playerPosition, cardPosition)
	models.NewAction(gameId, action)
	models.UpdateGameState(gameId, state)
	this.SetSuccessResponse()
}

func (this *ApiGameController) GameInfoCardValue() {
	gameId, _ := this.GetInt("game_id")
	state, err := models.ReadCurrentGameState(gameId)
	if this.SetError(err) {
		return
	}

	playerPosition, _ := this.GetInt("player_position")
	cardValue, _ := this.GetInt("card_value")
	action, _ := state.NewActionInformationValue(playerPosition, gamePackage.CardValue(cardValue))
	models.NewAction(gameId, action)
	models.UpdateGameState(gameId, state)
	this.SetSuccessResponse()
}

func (this *ApiGameController) GameInfoCardColor() {
	gameId, _ := this.GetInt("game_id")
	state, err := models.ReadCurrentGameState(gameId)
	if this.SetError(err) {
		return
	}

	playerPosition, _ := this.GetInt("player_position")
	cardColor, _ := this.GetInt("card_color")
	action, _ := state.NewActionInformationColor(playerPosition, gamePackage.CardColor(cardColor))
	models.NewAction(gameId, action)
	models.UpdateGameState(gameId, state)
	this.SetSuccessResponse()
}

func (this *ApiGameController) GameCurrentStep() {
	gameId, _ := this.GetInt("game_id")
	count, err := models.GetActionCount(gameId)
	if this.SetError(err) {
		return
	}
	result := struct {
		Status string `json:"status"`
		Step   int    `json:"step"`
	}{StatusSuccess, count}
	this.SetData(&result)
}

func (this *ApiGameController) GameInfo() {
	gameId, _ := this.GetInt("game_id")
	state, err := models.ReadCurrentGameState(gameId)
	if this.SetError(err) {
		return
	}

	playerPosition, err := state.GetPlayerPositionById(auth.GetUserIdFromSession(this.Ctx.Input.CruSession))
	if this.SetError(err) {
		return
	}

	result := struct {
		Status         string `json:"status"`
		Count          int    `json:"player_count"`
		PlayerPosition int    `json:"player_position"`
	}{StatusSuccess, state.PlayerCount, playerPosition}
	this.SetData(&result)
}
