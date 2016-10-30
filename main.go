package main

import (
	d "github.com/BabichMikhail/Hanabi/database"
	m "github.com/BabichMikhail/Hanabi/middlewares"
	_ "github.com/BabichMikhail/Hanabi/routers"
	"github.com/astaxie/beego"
	_ "github.com/mattn/go-sqlite3"
)

//var globalSessions *session.Manager

func init() {
	d.InitDatabase()
	m.InitMiddleware()
}

func main() {
	beego.Run()
}
