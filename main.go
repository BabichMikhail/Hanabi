package main

import (
	"fmt"

	gamePackage "github.com/BabichMikhail/Hanabi/game"
	"github.com/BabichMikhail/Hanabi/middlewares"
	"github.com/BabichMikhail/Hanabi/models"
	_ "github.com/BabichMikhail/Hanabi/routers"
	"github.com/Unknwon/goconfig"
	"github.com/astaxie/beego"
	"github.com/beego/compress"
	"github.com/beego/wetalk/setting"
)

func SettingCompress() {
	isProductMode := true
	setting, err := compress.LoadJsonConf("conf/compress.json", isProductMode, "")
	if err != nil {
		beego.Error(err)
		return
	}

	setting.RunCommand()
	if isProductMode {
		setting.RunCompress(true, false, true)
	}
	beego.AddFuncMap("compress_js", setting.Js.CompressJs)
	beego.AddFuncMap("compress_css", setting.Css.CompressCss)
}

func decrease(value interface{}) int {
	return value.(int) - 1
}

func init() {
	beego.SetStaticPath("/images", "static/images")
	beego.SetStaticPath("/css", "static/css")
	beego.SetStaticPath("/js", "static/js")

	beego.AddFuncMap("decrease", decrease)

	gamePackage.RegisterFunction()
	middlewares.InitMiddleware()

	Cfg, err := goconfig.LoadConfigFile("conf/app.conf")
	if err != nil {
		panic(err)
	}
	Cfg.BlockMode = false

	beego.BConfig.WebConfig.EnableXSRF = false

	driverName := Cfg.MustValue("orm", "driver_name", "mysql")
	dataSource := Cfg.MustValue("orm", "data_source", "root:root@/hanabi?charset=utf8")
	maxIdle := Cfg.MustInt("orm", "max_idle_conn", 30)
	maxOpen := Cfg.MustInt("orm", "max_open_conn", 50)

	models.RegisterDatabase(driverName, dataSource, maxIdle, maxOpen)

	setting.SecretKey = Cfg.MustValue("app", "secret_key")
	if len(setting.SecretKey) == 0 {
		fmt.Println("Please set your secret_key in app.conf file")
	}

	setting.LoginRememberDays = Cfg.MustInt("app", "login_remember_days", 7)
	setting.LoginMaxRetries = Cfg.MustInt("app", "login_max_retries", 5)
	setting.LoginFailedBlocks = Cfg.MustInt("app", "login_failed_blocks", 10)

	setting.CookieRememberName = Cfg.MustValue("app", "cookie_remember_name", "hanabi_remember_name")
	setting.CookieUserName = Cfg.MustValue("app", "cookie_user_name", "hanabi_cookie_name")
}

func main() {
	SettingCompress()
	beego.Run()
}
