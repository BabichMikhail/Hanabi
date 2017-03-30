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
	beego.Router("/admin/games/create/discardusefull", &controllers.AdminController{}, "get:GameDiscardUsefullCreate")
	beego.Router("/admin/games/create/usefullinfo", &controllers.AdminController{}, "get:GameUsefullInformationCreate")
	beego.Router("/admin/games/points/update", &controllers.AdminController{}, "get:UpdatePoints")

	beego.Router("/admin/stat/games/usefullinfo/:count_games/:count_players", &controllers.AdminController{}, "get:GameUsefullInformationRun")

	beego.Router("/admin", &controllers.AdminController{}, "get:Home")
	beego.Router("/admin/coefs", &controllers.AdminController{}, "get:FindUsefulInformationCoefs")
}
