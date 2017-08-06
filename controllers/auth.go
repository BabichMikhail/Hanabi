package controllers

import (
	"errors"
	"fmt"
	"regexp"

	"github.com/beego/wetalk/modules/auth"
	wetalk "github.com/beego/wetalk/modules/models"
)

type AuthController struct {
	BaseController
}

func (c *AuthController) SignIn() {
	c.SetBaseLayout()
	c.TplName = "templates/signin.html"
	if !c.Ctx.Input.IsPost() {
		return
	}
	username := c.GetString("username")
	password := c.GetString("password")
	var user wetalk.User
	if auth.VerifyUser(&user, username, password) {
		auth.LoginUser(&user, c.Ctx, true)
		c.Ctx.Redirect(302, c.URLFor("MainController.Get"))
	} else {
		c.Data["err"] = errors.New("Invalid credentials")
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

func (c *AuthController) SignUp() {
	c.SetBaseLayout()
	c.TplName = "templates/signup.html"
	if !c.Ctx.Input.IsPost() {
		return
	}
	username := c.GetString("username")
	if err := checkReservedUsernames(username); err != nil {
		c.Data["err"] = err
		return
	}
	email := c.GetString("email")
	password1 := c.GetString("password")
	password2 := c.GetString("confirmpassword")
	if password1 != password2 {
		c.Data["err"] = "Passwords don't match"
		return
	}

	var user wetalk.User
	if count, err := wetalk.Users().Count(); err == nil && count == 0 {
		user.IsAdmin = true
	} else if err != nil {
		panic(err)
	}
	err := auth.RegisterUser(&user, username, email, password1)

	if err == nil {
		c.Ctx.Redirect(302, c.URLFor(".SignIn"))
		return
	} else {
		c.Data["err"] = err
	}
}

func (c *AuthController) SignOut() {
	auth.LogoutUser(c.Ctx)
	c.Ctx.Redirect(302, c.URLFor(".SignIn"))
}
