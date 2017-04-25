package controllers

import (
	"strconv"

	gamePackage "github.com/BabichMikhail/Hanabi/game"
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
	gameInfo := state.GetPlayerGameInfo(userId, gamePackage.InfoTypeUsually)

	c.SetBaseLayout()
	c.TplName = "templates/game.html"

	var user wetalk.User
	auth.GetUserFromSession(&user, c.Ctx.Input.CruSession)
	c.Data["user"] = user
	c.LayoutSections = make(map[string]string)
	c.LayoutSections["Header"] = "components/navbar.html"
	c.LayoutSections["Scripts"] = "scripts/gamescripts.tpl"

	nickNames := make([]string, len(state.PlayerStates), len(state.PlayerStates))
	for i := 0; i < len(state.PlayerStates); i++ {
		nickNames[i] = models.GetUserNickNameById(state.PlayerStates[i].PlayerId)
	}

	var urls []CardUrl
	for _, color := range gamePackage.Colors {
		for _, value := range gamePackage.Values {
			urls = append(urls, CardUrl{
				Color: color,
				Value: value,
				Url:   gamePackage.GetCardUrlByValueAndColor(color, value),
			})
		}
	}

	c.Data["CardUrls"] = urls
	c.Data["PlayerInfo"] = gameInfo
	c.Data["NickNames"] = nickNames
	c.Data["Step"], _ = models.GetActionCount(id)
	c.Data["Players"] = models.GetGamePlayers([]int{id})[id]
	c.Data["MaxRedTokens"] = gamePackage.MaxRedTokens
	c.Data["MaxBlueTokens"] = gamePackage.MaxBlueTokens
	c.Data["NoneColor"] = gamePackage.NoneColor
	c.Data["NoneValue"] = gamePackage.NoneValue
	c.Data["TableColors"] = gamePackage.GetTableColorOrder()
	c.Data["ActionTypes"] = map[string]int{
		"infoColor": gamePackage.TypeActionInformationColor,
		"infoValue": gamePackage.TypeActionInformationValue,
		"discard":   gamePackage.TypeActionDiscard,
		"play":      gamePackage.TypeActionPlaying,
	}
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
