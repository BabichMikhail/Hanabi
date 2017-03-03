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

	beego.Router("/games/room/:id", &controllers.GameController{}, "get,post:Game")
	beego.Router("/games/finished/:id", &controllers.GameController{}, "get,post:GameFinished")

	beego.Router("/games/view/:id", &controllers.GameViewController{}, "get:GameView")

	beego.Router("/admin/games/create/random", &controllers.AdminController{}, "get:GameRandomCreate")
	beego.Router("/admin/games/create/smartyrandom", &controllers.AdminController{}, "get:GameSmartyRandomCreate")
}
