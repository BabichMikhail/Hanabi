package main

import (
	d "github.com/BabichMikhail/Hanabi/database"
	"github.com/BabichMikhail/Hanabi/engine"
	m "github.com/BabichMikhail/Hanabi/middlewares"
	_ "github.com/BabichMikhail/Hanabi/routers"
	"github.com/astaxie/beego"
	"github.com/beego/compress"
	_ "github.com/mattn/go-sqlite3"
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
	d.InitDatabase()
	beego.AddFuncMap("decrease", decrease)
	beego.AddFuncMap("cardValue", engine.GetCardValue)
	beego.AddFuncMap("cardColor", engine.GetCardColor)
	m.InitMiddleware()
}

func main() {
	SettingCompress()
	beego.Run()
}
