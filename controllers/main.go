package controllers

type MainController struct {
	BaseController
}

func (c *MainController) Get() {
	c.Ctx.Redirect(302, "/games")
	c.Data["Website"] = "beego.me"
	c.Data["Email"] = "astaxie@gmail.com"
	c.TplName = "index.tpl"
}
