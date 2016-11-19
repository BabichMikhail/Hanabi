package main

import (
	d "github.com/BabichMikhail/Hanabi/database"
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

func init() {
	beego.SetStaticPath("/images", "static/images")
	beego.SetStaticPath("/css", "static/css")
	beego.SetStaticPath("/js", "static/js")
	d.InitDatabase()
	m.InitMiddleware()
}

func main() {
	SettingCompress()
	beego.Run()
}
