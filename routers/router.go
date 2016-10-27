package routers

import (
	"github.com/BabichMikhail/Hanabi/controllers"
	"github.com/astaxie/beego"
)

func init() {
    beego.Router("/", &controllers.MainController{})
}
