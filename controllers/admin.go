package controllers

import (
	"math"

	ai "github.com/BabichMikhail/Hanabi/AI"
	"github.com/BabichMikhail/Hanabi/models"
	stat "github.com/BabichMikhail/Hanabi/statistics"
	"github.com/beego/wetalk/modules/auth"
)

type AdminController struct {
	BaseController
}

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
	c.GameCreate(ai.Type_AIRandom)
}

func (c *AdminController) GameSmartyRandomCreate() {
	c.GameCreate(ai.Type_AISmartyRandom)
}

func (c *AdminController) GameDiscardUsefullCreate() {
	c.GameCreate(ai.Type_AIDiscardUsefulCard)
}

func (c *AdminController) GameUsefullInformationCreate() {
	c.GameCreate(ai.Type_AIUsefulInformation)
}

func (c *AdminController) GameUsefulAndMaxMaxCreate() {
	c.GameCreate(ai.Type_AIUsefulInfoAndMaxMax)
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
		ai.Type_AIUsefulInformation,
		ai.Type_AIUsefulInformation,
		ai.Type_AIUsefulInformation,
		ai.Type_AIUsefulInformation,
		ai.Type_AIUsefulInformation,
	}
	stat.RunGames(aiTypes[0:countPlayers-1], []int{1, 2, 3, 4, 5}, countGames)
	c.Ctx.Redirect(302, c.URLFor("LobbyController.GameList"))
}

func (c *AdminController) FindUsefulInformationCoefs() {
	go stat.FindUsefulInfoV2Coefs()
	c.Ctx.Redirect(302, c.URLFor("AdminController.Home"))
}

func (c *AdminController) FindUsefulInformationCoefsV3() {
	gen := stat.NewGeneticAlgorithm()
	go gen.FindUsefulInfoV3Coefs()
	c.Ctx.Redirect(302, c.URLFor("AdminController.Home"))
}
