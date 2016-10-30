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
}
