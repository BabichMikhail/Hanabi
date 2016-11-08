package middlewares

import (
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

	ctx.Redirect(302, "/signin")
}

func InitMiddleware() {
	beego.InsertFilter("/", beego.BeforeRouter, CheckAuth)
	beego.InsertFilter("/games", beego.BeforeRouter, CheckAuth)
}
