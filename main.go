package main

import (
	_ "gatssh/routers"
	"github.com/astaxie/beego"
)

func main() {
	beego.Run()
}

func init() {
	beego.BConfig.WebConfig.Session.SessionOn = true
	beego.BConfig.WebConfig.Session.SessionName = "GatsshSessionID"
}
