package controllers

import (
	gamePackage "github.com/BabichMikhail/Hanabi/game"
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
	models.ApplyAction(gameId, gamePackage.TypeActionPlaying, playerPosition, cardPosition)
	this.SetSuccessResponse()
	go models.CheckAI(gameId)
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
	models.ApplyAction(gameId, gamePackage.TypeActionDiscard, playerPosition, cardPosition)
	this.SetSuccessResponse()
	go models.CheckAI(gameId)
}

func (this *ApiGameController) GameInfoCardValue() {
	gameId, _ := this.GetInt("game_id")
	playerPosition, _ := this.GetInt("player_position")
	cardValue, _ := this.GetInt("card_value")

	err := models.ApplyAction(gameId, gamePackage.TypeActionInformationValue, playerPosition, cardValue)
	if this.SetError(err) {
		return
	}

	this.SetSuccessResponse()
	go models.CheckAI(gameId)
}

func (this *ApiGameController) GameInfoCardColor() {
	gameId, _ := this.GetInt("game_id")
	playerPosition, _ := this.GetInt("player_position")
	cardColor, _ := this.GetInt("card_color")

	err := models.ApplyAction(gameId, gamePackage.TypeActionInformationColor, playerPosition, cardColor)
	if this.SetError(err) {
		return
	}

	this.SetSuccessResponse()
	go models.CheckAI(gameId)
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
	}{StatusSuccess, len(state.PlayerStates), playerPosition}
	this.SetData(&result)
}
