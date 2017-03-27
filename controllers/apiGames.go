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

func (c *ApiGameController) GamePlayCard() {
	gameId, _ := c.GetInt("game_id")
	state, err := models.ReadCurrentGameState(gameId)
	if c.SetFail(err) {
		return
	}

	playerPosition, err := state.GetPlayerPositionById(auth.GetUserIdFromSession(c.Ctx.Input.CruSession))
	if c.SetFail(err) {
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
	if c.SetFail(err) {
		return
	}

	playerPosition, err := state.GetPlayerPositionById(auth.GetUserIdFromSession(c.Ctx.Input.CruSession))
	if c.SetFail(err) {
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
	if c.SetFail(err) {
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
	if c.SetFail(err) {
		return
	}

	c.SetSuccessResponse()
	go models.CheckAI(gameId)
}

func (c *ApiGameController) GameCurrentStep() {
	gameId, _ := c.GetInt("game_id")
	count, err := models.GetActionCount(gameId)
	if c.SetFail(err) {
		return
	}
	c.SetData(&count)
}

func (c *ApiGameController) GameInfo() {
	gameId, _ := c.GetInt("game_id")
	state, err := models.ReadCurrentGameState(gameId)
	if c.SetFail(err) {
		return
	}

	playerPosition, err := state.GetPlayerPositionById(auth.GetUserIdFromSession(c.Ctx.Input.CruSession))
	gameInfo := state.GetPlayerGameInfoByPos(playerPosition)
	if c.SetFail(err) {
		return
	}

	c.SetData(&gameInfo)
}
