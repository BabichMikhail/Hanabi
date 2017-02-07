package middlewares

import (
	"strconv"

	"github.com/BabichMikhail/Hanabi/models"
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/context"
	"github.com/beego/wetalk/modules/auth"
	wetalk "github.com/beego/wetalk/modules/models"
)

func CheckAuth(ctx *context.Context) {
	if id := auth.GetUserIdFromSession(ctx.Input.CruSession); id > 0 {
		return
	}

	var user wetalk.User
	if auth.LoginUserFromRememberCookie(&user, ctx) {
		return
	}

	ctx.Redirect(401, "/signin")
}

func CheckUserInGame(ctx *context.Context) {
	userId := auth.GetUserIdFromSession(ctx.Input.CruSession)
	gameId, _ := strconv.Atoi(ctx.Input.Param(":id"))
	players := models.GetGamePlayers([]int{gameId})[gameId]
	for i := 0; i < len(players); i++ {
		if players[i].Id == userId {
			return
		}
	}
	ctx.Redirect(302, "/games")
}

func InitMiddleware() {
	beego.InsertFilter("/", beego.BeforeRouter, CheckAuth)
	beego.InsertFilter("/games", beego.BeforeRouter, CheckAuth)
	beego.InsertFilter("/api/*", beego.BeforeRouter, CheckAuth)
	beego.InsertFilter("/games/*", beego.BeforeRouter, CheckAuth)
	beego.InsertFilter("/games/room/:id", beego.BeforeRouter, CheckUserInGame)
}
