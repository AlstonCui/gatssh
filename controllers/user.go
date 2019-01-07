package controllers

import (
	"gatlin/models"
	"gatlin/utils"
	"fmt"
)

type UserLogin struct {
	baseController
}

func (this *UserLogin) Get() {
	this.TplName = "login.html"
}

type UserController struct {
	baseController
}

func (this *UserController) LoginAuth() {
	username := this.GetString("username")
	password := utils.Md5Sum(this.GetString("password"))

	user := &models.User{Username: username, Password: password,}

	clientIp:=this.getClientIp()

	Auth, uid := user.AuthUserAndPass()
	if Auth {
		sessionId:= fmt.Sprint(uid+clientIp)
		this.SetSession(GAT_SESSION_KEY, sessionId)
		this.SetSession(GAT_SESSION_USER,username)
		this.Redirect("/", 302)

		return
	} else {
		this.Redirect("/login", 302)
	}

}