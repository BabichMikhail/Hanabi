package controllers

import (
	"math"

	ai "github.com/BabichMikhail/Hanabi/AI"
	"github.com/BabichMikhail/Hanabi/models"
	stat "github.com/BabichMikhail/Hanabi/statistic"
	"github.com/beego/wetalk/modules/auth"
)

type AdminController struct {
	BaseController
}

// @todo api admin for create, read stats
func (c *AdminController) Home() {
	c.SetBaseLayout()
	c.TplName = "templates/adminhome.html"
	c.LayoutSections = make(map[string]string)
	c.LayoutSections["Header"] = "components/navbar.html"
	c.LayoutSections["Scripts"] = "scripts/adminscripts.tpl"
	c.Data["Stats"] = models.ReadStats()
}

func (c *AdminController) UpdatePoints() {
	userId := auth.GetUserIdFromSession(c.Ctx.Input.CruSession)
	games := models.GetFinishedGames(userId)
	for _, game := range games {
		models.UpdatePoints(game.Id)
	}

	c.Ctx.Redirect(302, c.URLFor("LobbyController.GameList"))
}

func (c *AdminController) GameCreate(aiType int) {
	count := 5
	userIds, err := models.GetAIUserIds(aiType, count)
	if err != nil {
		userIds, err = models.CreateAIUsers(aiType)
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

func (c *AdminController) GameRandomCreate() {
	c.GameCreate(ai.AI_RandomAction)
}

func (c *AdminController) GameSmartyRandomCreate() {
	c.GameCreate(ai.AI_SmartyRandomAction)
}

func (c *AdminController) GameDiscardUsefullCreate() {
	c.GameCreate(ai.AI_DiscardUsefullCardAction)
}

func (c *AdminController) GameUsefullInformationCreate() {
	c.GameCreate(ai.AI_UsefullInformationAction)
}

func (c *AdminController) GameUsefullInformationRun() {
	countGames, err := c.GetInt(":count_games")
	if err != nil {
		panic(err)
	}
	countPlayers, err := c.GetInt(":count_players")
	if err != nil {
		panic(err)
	}
	aiTypes := []int{
		ai.AI_UsefullInformationAction,
		ai.AI_UsefullInformationAction,
		ai.AI_UsefullInformationAction,
		ai.AI_UsefullInformationAction,
		ai.AI_UsefullInformationAction,
	}
	stat.RunGames(aiTypes, countGames, countPlayers)
	c.Ctx.Redirect(302, c.URLFor("LobbyController.GameList"))
}
