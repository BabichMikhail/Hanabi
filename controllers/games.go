package controllers

import (
	"fmt"
	"strconv"

	engine "github.com/BabichMikhail/Hanabi/engine"
	engineModels "github.com/BabichMikhail/Hanabi/engine/models"
	"github.com/BabichMikhail/Hanabi/models"
	"github.com/beego/wetalk/modules/auth"
	wetalk "github.com/beego/wetalk/modules/models"
)

type GameController struct {
	BaseController
}

func (this *GameController) GameList() {
	this.Layout = "base.tpl"
	this.TplName = "templates/gamelist.html"
	var user wetalk.User
	auth.GetUserFromSession(&user, this.Ctx.Input.CruSession)
	this.Data["user"] = user
	this.Data["games"] = models.GetGameList(engine.GetAllStatuses(), user.Id)
	if !this.Ctx.Input.IsPost() {
		return
	}
}

func (this *GameController) GameCreate() {
	var user wetalk.User
	playersCount, _ := this.GetInt("playersCount")
	auth.GetUserFromSession(&user, this.Ctx.Input.CruSession)
	id := models.NewGame(user.Id, playersCount, engineModels.GameWait)
	game, err := engine.MakeGame(id, user)
	if err != nil {
		result := struct {
			Status string `json:"status"`
			Err    error  `json:"err"`
		}{"FAIL", err}
		this.Data["json"] = &result
	} else {
		userId := auth.GetUserIdFromSession(this.Ctx.Input.CruSession)
		result := struct {
			Status      string            `json:"status"`
			Game        engineModels.Game `json:"game"`
			userId      int               `json:"currentUserId"`
			Err         error             `json:"err"`
			RedirectURL string            `json:"redirectURL"`
		}{"OK", game, userId, nil, this.URLFor("GameController.Game", ":id", id)}
		this.Data["json"] = &result
	}
	this.ServeJSON()
}

func (this *GameController) GameJoin() {
	id, _ := strconv.Atoi(this.Ctx.Input.Param(":id"))
	userId := auth.GetUserIdFromSession(this.Ctx.Input.CruSession)
	err, game_status := models.JoinGame(id, userId)
	var status string
	if err == nil {
		status = "OK"
	} else {
		status = "FAIL"
	}
	result := struct {
		Status      string `json:"status"`
		GameStatus  string `json:"game_status"`
		GameRoomURL string `json:"URL"`
	}{status, game_status, this.URLFor(".Game", ":id", id)}
	this.Data["json"] = &result
	this.ServeJSON()
}

func (this *GameController) GameLeave() {
	id, _ := strconv.Atoi(this.Ctx.Input.Param(":id"))
	userId := auth.GetUserIdFromSession(this.Ctx.Input.CruSession)
	action, err := models.LeaveGame(id, userId)
	var status string
	if err == nil && action != "" {
		status = "OK"
	} else {
		status = "FAIL"
	}
	result := struct {
		Status string `json:"status"`
		Action string `json:"action"`
	}{status, action}
	this.Data["json"] = &result
	this.ServeJSON()
}

func ConvertStringArrayToIntArray(s []string) []int {
	ans := []int{}
	for _, v := range s {
		i, _ := strconv.Atoi(v)
		ans = append(ans, i)
	}
	return ans
}

func (this *GameController) GameUpdate() {
	userId := auth.GetUserIdFromSession(this.Ctx.Input.CruSession)
	GameStatuses := models.GetStatuses(userId)
	for i, g := range GameStatuses {
		GameStatuses[i].URL = this.URLFor("GameController.Game", ":id", g.GameId)
	}
	result := struct {
		GameStatuses []models.GameStatus `json:"game"`
		Status       string              `json:"status"`
	}{GameStatuses, "OK"}
	this.Data["json"] = &result
	this.ServeJSON()
}

func (this *GameController) GameUsers() {
	id, _ := strconv.Atoi(this.Ctx.Input.Param(":id"))
	result := struct {
		Status  string                `json:"status"`
		Players []engineModels.Player `json:"players"`
	}{"OK", models.GetGamePlayers([]int{id})[id]}
	this.Data["json"] = &result
	this.ServeJSON()
}

func (this *GameController) Game() {
	id, _ := strconv.Atoi(this.Ctx.Input.Param(":id"))
	gamePlayers := models.GetGamePlayers([]int{id})[id]
	playerIds := []int{}
	for i := 0; i < len(gamePlayers); i++ {
		playerIds = append(playerIds, gamePlayers[i].Id)
	}
	game := engine.NewGame(playerIds)
	fmt.Println(game)
	this.Data["Website"] = "beego.me"
	this.Data["Email"] = "astaxie@gmail.com"
	this.TplName = "index.tpl"
}
