package controllers

import (
	"fmt"
	"strconv"

	engine "github.com/BabichMikhail/Hanabi/engine"
	"github.com/BabichMikhail/Hanabi/models"
)

type GameController struct {
	BaseController
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
