package controllers

import (
	"strconv"

	"github.com/BabichMikhail/Hanabi/models"
	"github.com/beego/wetalk/modules/auth"
	wetalk "github.com/beego/wetalk/modules/models"
)

type ApiLobbyController struct {
	ApiController
}

func (this *ApiLobbyController) GameCreate() {
	var user wetalk.User
	playersCount, _ := this.GetInt("playersCount")
	auth.GetUserFromSession(&user, this.Ctx.Input.CruSession)
	game, err := models.NewGame(user.Id, playersCount, models.StatusWait, false)
	if this.SetError(err) {
		return
	}

	userId := auth.GetUserIdFromSession(this.Ctx.Input.CruSession)
	data := struct {
		Id          int                  `json:"id"`
		Owner       string               `json:"owner"`
		Status      string               `json:"status"`
		Players     []models.LobbyPlayer `json:"players"`
		UserId      int                  `json:"currentUserId"`
		RedirectURL string               `json:"redirectURL"`
	}{game.Id, game.Owner, game.Status, game.Players, userId, this.URLFor("GameController.Game", ":id", game.Id)}
	result := struct {
		Status string      `json:"status"`
		Data   interface{} `json:"data"`
	}{StatusSuccess, data}
	this.SetData(&result)
}

func (this *ApiLobbyController) GameJoin() {
	id, _ := strconv.Atoi(this.Ctx.Input.Param(":id"))
	userId := auth.GetUserIdFromSession(this.Ctx.Input.CruSession)
	err, gameStatus := models.JoinGame(id, userId)
	if this.SetError(err) {
		return
	}

	result := struct {
		Status      string `json:"status"`
		GameStatus  string `json:"game_status"`
		GameRoomURL string `json:"URL"`
	}{StatusSuccess, gameStatus, this.URLFor(".Game", ":id", id)}
	this.SetData(&result)
}

func (this *ApiLobbyController) GameLeave() {
	id, _ := strconv.Atoi(this.Ctx.Input.Param(":id"))
	userId := auth.GetUserIdFromSession(this.Ctx.Input.CruSession)
	action, err := models.LeaveGame(id, userId)
	if this.SetError(err) {
		return
	}

	result := struct {
		Status string `json:"status"`
		Action string `json:"action"`
	}{StatusSuccess, action}
	this.SetData(&result)
}

func ConvertStringArrayToIntArray(s []string) []int {
	ans := []int{}
	for _, v := range s {
		i, _ := strconv.Atoi(v)
		ans = append(ans, i)
	}
	return ans
}

func (this *ApiLobbyController) setGameURLs(games []models.LobbyGame) {
	for i, g := range games {
		games[i].URL = this.URLFor("GameController.Game", ":id", g.Id)
	}
}

func (this *ApiLobbyController) setGameData(games []models.LobbyGame) {
	result := struct {
		Status string             `json:"status"`
		Games  []models.LobbyGame `json:"games"`
	}{StatusSuccess, games}
	this.SetData(&result)
}

func (c *ApiLobbyController) setGames(getGames func(int) []models.LobbyGame) {
	userId := auth.GetUserIdFromSession(c.Ctx.Input.CruSession)
	games := getGames(userId)
	c.setGameURLs(games)
	c.setGameData(games)
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

func (this *ApiLobbyController) MyInfo() {
	var user wetalk.User
	auth.GetUserFromSession(&user, this.Ctx.Input.CruSession)
	userResult := UserInfo{user.Id, user.NickName}
	result := struct {
		Status string   `json:"status"`
		User   UserInfo `json:"user"`
	}{StatusSuccess, userResult}
	this.SetData(&result)
}
