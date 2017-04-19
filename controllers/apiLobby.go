package controllers

import (
	"strconv"

	ai "github.com/BabichMikhail/Hanabi/AI"
	"github.com/BabichMikhail/Hanabi/models"
	"github.com/beego/wetalk/modules/auth"
	wetalk "github.com/beego/wetalk/modules/models"
)

type ApiLobbyController struct {
	ApiController
}

func (c *ApiLobbyController) GameCreate() {
	var user wetalk.User
	playersCount, _ := c.GetInt("playersCount")
	auth.GetUserFromSession(&user, c.Ctx.Input.CruSession)
	game, err := models.NewGame(user.Id, playersCount, models.StatusWait, false)
	if c.SetFail(err) {
		return
	}

	data := struct {
		Id          int                  `json:"id"`
		Owner       string               `json:"owner"`
		Status      string               `json:"status"`
		Players     []models.LobbyPlayer `json:"players"`
		UserId      int                  `json:"currentUserId"`
		RedirectURL string               `json:"redirectURL"`
	}{game.Id, game.Owner, game.Status, game.Players, user.Id, c.URLFor("GameController.Game", ":id", game.Id)}
	c.SetData(&data)
}

func (c *ApiLobbyController) GameAICreate() {
	var user wetalk.User
	playersCount, _ := c.GetInt("playersCount")
	auth.GetUserFromSession(&user, c.Ctx.Input.CruSession)
	game, err := models.NewGame(user.Id, playersCount, models.StatusWait, false)
	var gameStatus string
	for i := 1; i < playersCount; i++ {
		var userId int
		userId, err = models.GetAIUserId(ai.Type_AIUsefulInfoV4AndParts, i)
		if c.SetFail(err) {
			return
		}

		err, gameStatus = models.JoinGame(game.Id, userId)
		if c.SetFail(err) {
			return
		}

	}

	if c.SetFail(err) {
		return
	}

	players := models.GetGamePlayers([]int{game.Id})[game.Id]
	data := struct {
		Id          int                  `json:"id"`
		Owner       string               `json:"owner"`
		Status      string               `json:"status"`
		Players     []models.LobbyPlayer `json:"players"`
		UserId      int                  `json:"currentUserId"`
		RedirectURL string               `json:"redirectURL"`
	}{game.Id, game.Owner, gameStatus, players, user.Id, c.URLFor("GameController.Game", ":id", game.Id)}
	c.SetData(&data)
	go models.CheckAI(game.Id)
}

func (c *ApiLobbyController) GameJoin() {
	id, _ := strconv.Atoi(c.Ctx.Input.Param(":id"))
	userId := auth.GetUserIdFromSession(c.Ctx.Input.CruSession)
	err, gameStatus := models.JoinGame(id, userId)
	if c.SetFail(err) {
		return
	}

	data := struct {
		GameStatus  string `json:"game_status"`
		GameRoomURL string `json:"URL"`
	}{gameStatus, c.URLFor(".Game", ":id", id)}
	c.SetData(&data)
}

func (c *ApiLobbyController) GameLeave() {
	id, _ := strconv.Atoi(c.Ctx.Input.Param(":id"))
	userId := auth.GetUserIdFromSession(c.Ctx.Input.CruSession)
	action, err := models.LeaveGame(id, userId)
	if c.SetFail(err) {
		return
	}

	c.SetData(&action)
}

func ConvertStringArrayToIntArray(s []string) []int {
	ans := []int{}
	for _, v := range s {
		i, _ := strconv.Atoi(v)
		ans = append(ans, i)
	}
	return ans
}

func (c *ApiLobbyController) setGameURLs(games []models.LobbyGame) {
	for i, g := range games {
		games[i].URL = c.URLFor("GameController.Game", ":id", g.Id)
	}
}

func (c *ApiLobbyController) setGames(getGames func(int) []models.LobbyGame) {
	userId := auth.GetUserIdFromSession(c.Ctx.Input.CruSession)
	games := getGames(userId)
	c.setGameURLs(games)
	c.SetData(&games)
}

func (c *ApiLobbyController) GetActiveGames() {
	c.setGames(models.GetActiveGames)
}

func (c *ApiLobbyController) GetMyGames() {
	c.setGames(models.GetMyGames)
}

func (c *ApiLobbyController) GetAllGames() {
	c.setGames(models.GetAllGames)
}

func (c *ApiLobbyController) GetFinishedGames() {
	c.setGames(models.GetFinishedGames)
}

type UserInfo struct {
	Id       int    `json:"id"`
	NickName string `json:"nick_name"`
}

func (c *ApiLobbyController) MyInfo() {
	var user wetalk.User
	auth.GetUserFromSession(&user, c.Ctx.Input.CruSession)
	userResult := UserInfo{user.Id, user.NickName}
	c.SetData(&userResult)
}
