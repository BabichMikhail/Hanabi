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

	beego.Router("/games", &controllers.GameController{}, "get,post:GameList")
	beego.Router("/games/room/:id", &controllers.GameController{}, "get,post:Game")

	beego.Router("/api/games/create", &controllers.GameController{}, "post:GameCreate")
	beego.Router("/api/games/status", &controllers.GameController{}, "get:GameUpdate")
	beego.Router("/api/games/join/:id", &controllers.GameController{}, "post:GameJoin")
	beego.Router("/api/games/leave/:id", &controllers.GameController{}, "post:GameLeave")
	beego.Router("/api/games/users/:id", &controllers.GameController{}, "get:GameUsers")

	beego.Router("/api/users/current", &controllers.AuthController{}, "get:UserCurrent")
}
