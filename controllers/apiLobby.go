package controllers

import (
	"strconv"

	engineLobby "github.com/BabichMikhail/Hanabi/engine/lobby"
	"github.com/BabichMikhail/Hanabi/models"
	"github.com/beego/wetalk/modules/auth"
	wetalk "github.com/beego/wetalk/modules/models"
)

type LobbyApiController struct {
	ApiController
}

func (this *LobbyApiController) GameCreate() {
	var user wetalk.User
	playersCount, _ := this.GetInt("playersCount")
	auth.GetUserFromSession(&user, this.Ctx.Input.CruSession)
	id := models.NewGame(user.Id, playersCount, engineLobby.GameWait)
	game, err := engineLobby.MakeGame(id, user)
	if this.SetError(err) {
		return
	}

	userId := auth.GetUserIdFromSession(this.Ctx.Input.CruSession)
	result := struct {
		Status      string           `json:"status"`
		Game        engineLobby.Game `json:"game"`
		UserId      int              `json:"currentUserId"`
		RedirectURL string           `json:"redirectURL"`
	}{StatusSuccess, game, userId, this.URLFor("GameController.Game", ":id", id)}
	this.SetData(&result)
}

func (this *LobbyApiController) GameJoin() {
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

func (this *LobbyApiController) GameLeave() {
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

func (this *LobbyApiController) GameUpdate() {
	userId := auth.GetUserIdFromSession(this.Ctx.Input.CruSession)
	GameStatuses := models.GetStatuses(userId)
	for i, g := range GameStatuses {
		GameStatuses[i].URL = this.URLFor("GameController.Game", ":id", g.Game.Id)
	}

	result := struct {
		Status       string              `json:"status"`
		GameStatuses []models.GameStatus `json:"games"`
	}{StatusSuccess, GameStatuses}
	this.SetData(&result)
}

func (this *LobbyApiController) GameUsers() {
	id, _ := strconv.Atoi(this.Ctx.Input.Param(":id"))
	result := struct {
		Status  string               `json:"status"`
		Players []engineLobby.Player `json:"players"`
	}{StatusSuccess, models.GetGamePlayers([]int{id})[id]}
	this.SetData(&result)
}

type UserInfo struct {
	Id       int    `json:"id"`
	NickName string `json:"nick_name"`
}

func (this *LobbyApiController) MyInfo() {
	var user wetalk.User
	auth.GetUserFromSession(&user, this.Ctx.Input.CruSession)
	userResult := UserInfo{user.Id, user.NickName}
	result := struct {
		Status string   `json:"status"`
		User   UserInfo `json:"user"`
	}{StatusSuccess, userResult}
	this.SetData(&result)
}
