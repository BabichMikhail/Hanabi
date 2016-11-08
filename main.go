package main

import (
	d "github.com/BabichMikhail/Hanabi/database"
	m "github.com/BabichMikhail/Hanabi/middlewares"
	_ "github.com/BabichMikhail/Hanabi/routers"
	"github.com/astaxie/beego"
	_ "github.com/mattn/go-sqlite3"
)

func init() {
	beego.SetStaticPath("/images", "static/images")
	beego.SetStaticPath("/css", "static/css")
	beego.SetStaticPath("/js", "static/js")
	d.InitDatabase()
	m.InitMiddleware()

}

func main() {
	beego.Run()
}
