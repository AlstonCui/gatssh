package controllers

import (
	"gatssh/models"
	"fmt"
)

type SettingController struct {
	baseController
}

func (this *SettingController) Setting() {

	if this.IsLogin == true {
		this.TplName = "setting.html"
		return
	}
	this.Redirect("/login", 302)
	return
}

func (this *SettingController) ChangePassword() {

	if this.IsLogin != true {
		this.ServeJSON(40000, "Please login...")
		return
	}

	oldPass := this.GetString("oldPass")
	newPass := this.GetString("newPass")
	confirmPass := this.GetString("confirmPass")

	if oldPass == "" || newPass == "" {
		this.ServeJSON(40000, "something is empty")
		return
	}

	if newPass != confirmPass {
		this.ServeJSON(30000, "confirm err")
		return
	}

	user := &models.User{Username: this.User, Password: oldPass,}

	fmt.Println(user)
	uid, Auth := user.AuthUserAndPass()

	if !Auth {
		this.ServeJSON(30000, "old password is wrong!")
		return
	}

	err := user.UpdatePassword(uid, newPass)
	if err != nil {
		this.ServeJSON(40000, err)
		return
	}

	this.ServeJSON(20000, nil)
	return

}
