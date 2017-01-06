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
	game, err := models.ReadActiveGameById(id)
	if err != nil {
		game, _ = models.CreateActiveGame(playerIds, id)
	}
	userId := auth.GetUserIdFromSession(this.Ctx.Input.CruSession)
	playerInfo := game.GetPlayerGameInfo(userId)
	this.Data["json"] = &playerInfo
	this.ServeJSON()
}

func (this *ApiGameController) GamePlayCard() {
	gameId, _ := this.GetInt("game_id")
	cardPosition, _ := this.GetInt("card_position")
	game, err := models.ReadActiveGameById(gameId)
	if err != nil {
		result := struct {
			Status string `json:"status"`
			Err    error  `json:"err"`
		}{"FAIL", err}
		this.Data["json"] = &result
		this.ServeJSON()
		return
	}
	playerPosition, err := game.GetPlayerPositionById(auth.GetUserIdFromSession(this.Ctx.Input.CruSession))
	if err != nil {
		result := struct {
			Status string `json:"status"`
			Err    error  `json:"err"`
		}{"FAIL", err}
		this.Data["json"] = &result
		this.ServeJSON()
		return
	}
	game.NewActionPlaying(playerPosition, cardPosition)
	result := struct {
		Status string `json:status`
	}{"OK"}
	this.Data["json"] = &result
	this.ServeJSON()
}

func (this *ApiGameController) GameDiscardCard() {
	gameId, _ := this.GetInt("game_id")
	cardPosition, _ := this.GetInt("card_position")
	game, err := models.ReadActiveGameById(gameId)
	if err != nil {
		result := struct {
			Status string `json:"status"`
			Err    error  `json:"err"`
		}{"FAIL", err}
		this.Data["json"] = &result
		this.ServeJSON()
		return
	}
	playerPosition, err := game.GetPlayerPositionById(auth.GetUserIdFromSession(this.Ctx.Input.CruSession))
	if err != nil {
		result := struct {
			Status string `json:"status"`
			Err    error  `json:"err"`
		}{"FAIL", err}
		this.Data["json"] = &result
		this.ServeJSON()
		return
	}
	game.NewActionDiscard(playerPosition, cardPosition)
	models.UpdateCurrentGameById(gameId, game)
	result := struct {
		Status string `json:status`
	}{"OK"}
	this.Data["json"] = &result
	this.ServeJSON()
}

func (this *ApiGameController) GameInfoCardValue() {
	gameId, _ := this.GetInt("game_id")
	playerPosition, _ := this.GetInt("player_position")
	cardValue, _ := this.GetInt("card_value")
	game, err := models.ReadActiveGameById(gameId)
	if err != nil {
		result := struct {
			Status string `json:"status"`
			Err    error  `json:"err"`
		}{"FAIL", err}
		this.Data["json"] = &result
		this.ServeJSON()
		return
	}
	game.NewActionInformationValue(playerPosition, gamePackage.CardValue(cardValue))
	models.UpdateCurrentGameById(gameId, game)
	result := struct {
		Status string `json:status`
	}{"OK"}
	this.Data["json"] = &result
	this.ServeJSON()
}

func (this *ApiGameController) GameInfoCardColor() {
	gameId, _ := this.GetInt("game_id")
	playerPosition, _ := this.GetInt("player_position")
	cardColor, _ := this.GetInt("card_color")
	game, err := models.ReadActiveGameById(gameId)
	if err != nil {
		result := struct {
			Status string `json:"status"`
			Err    error  `json:"err"`
		}{"FAIL", err}
		this.Data["json"] = &result
		this.ServeJSON()
		return
	}
	game.NewActionInformationColor(playerPosition, gamePackage.CardColor(cardColor))
	models.UpdateCurrentGameById(gameId, game)
	result := struct {
		Status string `json:status`
	}{"OK"}
	this.Data["json"] = &result
	this.ServeJSON()
}
