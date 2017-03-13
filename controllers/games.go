package controllers

import (
	"strconv"

	engineGame "github.com/BabichMikhail/Hanabi/game"
	"github.com/BabichMikhail/Hanabi/models"
	"github.com/beego/wetalk/modules/auth"
	wetalk "github.com/beego/wetalk/modules/models"
)

type GameController struct {
	BaseController
}

func (c *GameController) Game() {
	id, _ := strconv.Atoi(c.Ctx.Input.Param(":id"))
	state, err := models.ReadCurrentGameState(id)

	if err != nil {
		c.Ctx.Redirect(302, c.URLFor("LobbyController.GameList"))
		return
	}

	if state.IsGameOver() {
		models.SetGameFinishedStatus(id)
		c.Ctx.Redirect(302, c.URLFor("GameController.GameFinished", ":id", id))
		return
	}

	userId := auth.GetUserIdFromSession(c.Ctx.Input.CruSession)
	playerInfo := state.GetPlayerGameInfo(userId)
	c.Data["playerInfo"] = playerInfo
	c.Data["deckFirstNumber"] = playerInfo.DeckSize / 10
	c.Data["deckSecondNumber"] = playerInfo.DeckSize % 10

	c.Data["Step"], _ = models.GetActionCount(id)
	c.SetBaseLayout()
	c.TplName = "templates/game.html"

	var user wetalk.User
	auth.GetUserFromSession(&user, c.Ctx.Input.CruSession)
	c.Data["user"] = user
	c.LayoutSections = make(map[string]string)
	c.LayoutSections["Header"] = "components/navbar.html"
	c.LayoutSections["Scripts"] = "scripts/gamescripts.tpl"
	c.Data["TableColors"] = engineGame.GetTableColorOrder()
	playerNickNames := []string{}
	for i := 0; i < len(state.PlayerStates); i++ {
		nickName := models.GetUserNickNameById(state.PlayerStates[i].PlayerId)
		playerNickNames = append(playerNickNames, nickName)
	}
	c.Data["NickNames"] = playerNickNames
}

func (c *GameController) GameFinished() {
	id, _ := strconv.Atoi(c.Ctx.Input.Param(":id"))

	state, err := models.ReadCurrentGameState(id)
	if err != nil {
		c.Ctx.Redirect(302, c.URLFor("LobbyController.GameList"))
		return
	}
	if !state.IsGameOver() {
		c.Ctx.Redirect(302, c.URLFor("GameController.Game", ":id", id))
		return
	}
	c.Data["Points"], _ = state.GetPoints()
	c.SetBaseLayout()
	c.TplName = "templates/gamefinished.html"

	var user wetalk.User
	auth.GetUserFromSession(&user, c.Ctx.Input.CruSession)
	c.Data["user"] = user
	c.LayoutSections = make(map[string]string)
	c.LayoutSections["Header"] = "components/navbar.html"
}
