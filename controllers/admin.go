package controllers

import (
	"math"

	ai "github.com/BabichMikhail/Hanabi/AI"
	"github.com/BabichMikhail/Hanabi/models"
)

type AdminController struct {
	BaseController
}

func (c *AdminController) GameCreate() {
	count := 5
	userIds, err := models.GetAIUserIds(ai.AI_RandomAction, count)
	if err != nil {
		userIds, err = models.CreateAIUsers(ai.AI_RandomAction)
		if err != nil {
			c.Ctx.Redirect(302, c.URLFor("LobbyController.GameList"))
		}
	}

	gameItem, _ := models.NewGame(userIds[0], count, models.StatusWait, true)
	gameId := gameItem.Id
	for i := 1; i < int(math.Min(float64(len(userIds)), float64(count))); i++ {
		models.JoinGame(gameId, userIds[i])
	}
	models.CheckAI(gameItem.Id)
	c.Ctx.Redirect(302, c.URLFor("ViewController.GameView", ":id", gameId))
}
