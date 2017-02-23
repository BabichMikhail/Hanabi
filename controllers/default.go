package controllers

import (
	"github.com/astaxie/beego"
)

type BaseController struct {
	beego.Controller
}

func (c *BaseController) SetBaseLayout() {
	c.Layout = "base.tpl"
}
