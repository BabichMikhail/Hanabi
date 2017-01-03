package routers

import (
	"github.com/BabichMikhail/Hanabi/controllers"
	"github.com/astaxie/beego"
)

func init() {
	beego.Router("/", &controllers.MainController{})

	beego.Router("/signin", &controllers.AuthController{}, "get,post:SignIn")
	beego.Router("/signup", &controllers.AuthController{}, "get,post:SignUp")
	beego.Router("/signout", &controllers.AuthController{}, "*:SignOut")

	beego.Router("/games", &controllers.LobbyController{}, "get,post:GameList")

	beego.Router("/api/games/create", &controllers.LobbyApiController{}, "post:GameCreate")
	beego.Router("/api/games/status", &controllers.LobbyApiController{}, "get:GameUpdate")
	beego.Router("/api/games/join/:id", &controllers.LobbyApiController{}, "post:GameJoin")
	beego.Router("/api/games/leave/:id", &controllers.LobbyApiController{}, "post:GameLeave")
	beego.Router("/api/games/users/:id", &controllers.LobbyApiController{}, "get:GameUsers")

	beego.Router("/api/users/current", &controllers.AuthController{}, "get:UserCurrent")

	beego.Router("/games/room/:id", &controllers.GameController{}, "get,post:Game")
}
