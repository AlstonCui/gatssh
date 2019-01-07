package main

import (
	_ "gatlin/routers"
	"github.com/astaxie/beego"
	"encoding/gob"
	"gatlin/models"
)



func main() {
	beego.Run()
}

func init() {
	//beego的session序列号是用gob的方式，因此需要将注册models.User
	gob.Register(models.User{})
	//https://beego.me/docs/mvc/controller/session.md
	beego.BConfig.WebConfig.Session.SessionOn = true
	beego.BConfig.WebConfig.Session.SessionName = "GatsshSessionID"
}
