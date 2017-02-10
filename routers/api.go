package routers

import (
	"github.com/BabichMikhail/Hanabi/controllers"
	"github.com/astaxie/beego"
)

func init() {
	beego.Router("/api/lobby/create", &controllers.ApiLobbyController{}, "post:GameCreate")
	beego.Router("/api/lobby/status", &controllers.ApiLobbyController{}, "get:GameUpdate")
	beego.Router("/api/lobby/join/:id", &controllers.ApiLobbyController{}, "post:GameJoin")
	beego.Router("/api/lobby/leave/:id", &controllers.ApiLobbyController{}, "post:GameLeave")
	beego.Router("/api/lobby/users/:id", &controllers.ApiLobbyController{}, "get:GameUsers")

	beego.Router("/api/games/cards", &controllers.ApiGameController{}, "get:GetGameCards")
	beego.Router("/api/games/statuses", &controllers.ApiGameController{}, "get:GetGameStatuses")

	beego.Router("/api/games/action/play", &controllers.ApiGameController{}, "post:GamePlayCard")
	beego.Router("/api/games/action/discard", &controllers.ApiGameController{}, "post:GameDiscardCard")
	beego.Router("/api/games/action/info/value", &controllers.ApiGameController{}, "post:GameInfoCardValue")
	beego.Router("/api/games/action/info/color", &controllers.ApiGameController{}, "post:GameInfoCardColor")
	beego.Router("/api/games/step", &controllers.ApiGameController{}, "get:GameCurrentStep")
	beego.Router("/api/games/info", &controllers.ApiGameController{}, "get:GameInfo")

	beego.Router("/api/users/current", &controllers.ApiLobbyController{}, "get:MyInfo")
}
