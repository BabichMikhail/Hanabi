package controllers

type MainController struct {
	BaseController
}

func (this *MainController) Get() {
	this.Ctx.Redirect(302, "/games")
	this.Data["Website"] = "beego.me"
	this.Data["Email"] = "astaxie@gmail.com"
	this.TplName = "index.tpl"
}
