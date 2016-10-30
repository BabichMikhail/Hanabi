package controllers

import (
	"errors"

	"github.com/beego/wetalk/modules/auth"
	wetalk "github.com/beego/wetalk/modules/models"
	_ "github.com/mattn/go-sqlite3"
)

type AuthController struct {
	BaseController
}

func (this *AuthController) SignIn() {
	this.Layout = "base.tpl"
	this.TplName = "templates/signin.html"
	if !this.Ctx.Input.IsPost() {
		return
	}
	username := this.GetString("username")
	password := this.GetString("password")
	var user wetalk.User
	if auth.VerifyUser(&user, username, password) {
		auth.LoginUser(&user, this.Ctx, true)
		this.Ctx.Redirect(302, this.URLFor("MainController.Get"))
	} else {
		this.Data["err"] = errors.New("Invalid credentials")
	}
}

func (this *AuthController) SignUp() {
	this.Layout = "base.tpl"
	this.TplName = "templates/signup.html"
	if !this.Ctx.Input.IsPost() {
		return
	}
	username := this.GetString("username")
	email := this.GetString("email")
	password1 := this.GetString("password")
	password2 := this.GetString("confirmpassword")
	if password1 != password2 {
		this.Data["err"] = "Password don't match"
		return
	}

	var user wetalk.User
	err := auth.RegisterUser(&user, username, email, password1)
	if err == nil {
		this.Ctx.Redirect(302, this.URLFor(".SignIn"))
	} else {
		this.Data["err"] = err
	}
}

func (this *AuthController) SignOut() {
	auth.LogoutUser(this.Ctx)
	this.Ctx.Redirect(302, this.URLFor(".SignIn"))
}
