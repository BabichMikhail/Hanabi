package routers

import (
	"github.com/BabichMikhail/Hanabi/controllers"
	"github.com/astaxie/beego"
)

func init() {
	NSApi := beego.NewNamespace("/api",
		// @todo add auth filter for api
		beego.NSNamespace("/lobby",
			beego.NSRouter("/create", &controllers.ApiLobbyController{}, "post:GameCreate"),
			beego.NSRouter("/join/:id", &controllers.ApiLobbyController{}, "post:GameJoin"),
			beego.NSRouter("/leave/:id", &controllers.ApiLobbyController{}, "post:GameLeave"),
			beego.NSNamespace("/games",
				beego.NSRouter("/active", &controllers.ApiLobbyController{}, "get:GetActiveGames"),
				beego.NSRouter("/my", &controllers.ApiLobbyController{}, "get:GetMyGames"),
				beego.NSRouter("/all", &controllers.ApiLobbyController{}, "get:GetAllGames"),
				beego.NSRouter("/finished", &controllers.ApiLobbyController{}, "get:GetFinishedGames"),
			),
		),

		beego.NSNamespace("/games",
			beego.NSRouter("/cards", &controllers.ApiGameController{}, "get:GetGameCards"),

			beego.NSNamespace("/action",
				beego.NSRouter("/play", &controllers.ApiGameController{}, "post:GamePlayCard"),
				beego.NSRouter("/discard", &controllers.ApiGameController{}, "post:GameDiscardCard"),
				beego.NSRouter("/info/value", &controllers.ApiGameController{}, "post:GameInfoCardValue"),
				beego.NSRouter("/info/color", &controllers.ApiGameController{}, "post:GameInfoCardColor"),
			),

			beego.NSRouter("/step", &controllers.ApiGameController{}, "get:GameCurrentStep"),
			beego.NSRouter("/info", &controllers.ApiGameController{}, "get:GameInfo"),
		),

		beego.NSNamespace("/ai",
			beego.NSRouter("/names", &controllers.ApiAdminController{}, "get:GetAINames"),
		),

		beego.NSNamespace("/users",
			beego.NSRouter("/current", &controllers.ApiLobbyController{}, "get:MyInfo"),
		),

		beego.NSNamespace("/admin",
			beego.NSNamespace("/stats",
				beego.NSRouter("/read", &controllers.ApiAdminController{}, "get:ReadStats"),
				beego.NSRouter("/create", &controllers.ApiAdminController{}, "post:CreateStat"),
			),
		),
	)

	beego.AddNamespace(NSApi)

}
