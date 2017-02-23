package controllers

import (
	"errors"
	"fmt"
	"regexp"

	"github.com/beego/wetalk/modules/auth"
	wetalk "github.com/beego/wetalk/modules/models"
	_ "github.com/mattn/go-sqlite3"
)

type AuthController struct {
	BaseController
}

func (this *AuthController) SignIn() {
	this.SetBaseLayout()
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

func checkReservedUsernames(username string) error {
	reservedPatterns := []string{
		"AI_.*",
	}
	for _, pattern := range reservedPatterns {
		r, _ := regexp.Compile(pattern)

		if r.MatchString(username) {
			return errors.New(fmt.Sprintf("Can't register user with username which matched to pattern: %s", pattern))
		}
	}
	return nil
}

func (this *AuthController) SignUp() {
	this.SetBaseLayout()
	this.TplName = "templates/signup.html"
	if !this.Ctx.Input.IsPost() {
		return
	}
	username := this.GetString("username")
	if err := checkReservedUsernames(username); err != nil {
		this.Data["err"] = err
		return
	}
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
		return
	} else {
		this.Data["err"] = err
	}
}

func (this *AuthController) SignOut() {
	auth.LogoutUser(this.Ctx)
	this.Ctx.Redirect(302, this.URLFor(".SignIn"))
}
