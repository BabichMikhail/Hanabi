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

func (c *ApiGameController) GetGameCards() {
	result := struct {
		Status string                           `json:"status"`
		Colors map[gamePackage.CardColor]string `json:"colors"`
		Values map[gamePackage.CardValue]string `json:"values"`
	}{StatusSuccess, card.GetColors(), card.GetValues()}
	c.SetData(&result)
}

func (c *ApiGameController) GamePlayCard() {
	gameId, _ := c.GetInt("game_id")
	state, err := models.ReadCurrentGameState(gameId)
	if c.SetError(err) {
		return
	}

	playerPosition, err := state.GetPlayerPositionById(auth.GetUserIdFromSession(c.Ctx.Input.CruSession))
	if c.SetError(err) {
		return
	}

	cardPosition, _ := c.GetInt("card_position")
	models.ApplyAction(gameId, gamePackage.TypeActionPlaying, playerPosition, cardPosition)
	c.SetSuccessResponse()
	go models.CheckAI(gameId)
}

func (c *ApiGameController) GameDiscardCard() {
	gameId, _ := c.GetInt("game_id")
	state, err := models.ReadCurrentGameState(gameId)
	if c.SetError(err) {
		return
	}

	playerPosition, err := state.GetPlayerPositionById(auth.GetUserIdFromSession(c.Ctx.Input.CruSession))
	if c.SetError(err) {
		return
	}

	cardPosition, _ := c.GetInt("card_position")
	models.ApplyAction(gameId, gamePackage.TypeActionDiscard, playerPosition, cardPosition)
	c.SetSuccessResponse()
	go models.CheckAI(gameId)
}

func (c *ApiGameController) GameInfoCardValue() {
	gameId, _ := c.GetInt("game_id")
	playerPosition, _ := c.GetInt("player_position")
	cardValue, _ := c.GetInt("card_value")

	err := models.ApplyAction(gameId, gamePackage.TypeActionInformationValue, playerPosition, cardValue)
	if c.SetError(err) {
		return
	}

	c.SetSuccessResponse()
	go models.CheckAI(gameId)
}

func (c *ApiGameController) GameInfoCardColor() {
	gameId, _ := c.GetInt("game_id")
	playerPosition, _ := c.GetInt("player_position")
	cardColor, _ := c.GetInt("card_color")

	err := models.ApplyAction(gameId, gamePackage.TypeActionInformationColor, playerPosition, cardColor)
	if c.SetError(err) {
		return
	}

	c.SetSuccessResponse()
	go models.CheckAI(gameId)
}

func (c *ApiGameController) GameCurrentStep() {
	gameId, _ := c.GetInt("game_id")
	count, err := models.GetActionCount(gameId)
	if c.SetError(err) {
		return
	}
	result := struct {
		Status string `json:"status"`
		Step   int    `json:"step"`
	}{StatusSuccess, count}
	c.SetData(&result)
}

func (c *ApiGameController) GameInfo() {
	gameId, _ := c.GetInt("game_id")
	state, err := models.ReadCurrentGameState(gameId)
	if c.SetError(err) {
		return
	}

	playerPosition, err := state.GetPlayerPositionById(auth.GetUserIdFromSession(c.Ctx.Input.CruSession))
	if c.SetError(err) {
		return
	}

	result := struct {
		Status         string `json:"status"`
		Count          int    `json:"player_count"`
		PlayerPosition int    `json:"player_position"`
	}{StatusSuccess, len(state.PlayerStates), playerPosition}
	c.SetData(&result)
}
