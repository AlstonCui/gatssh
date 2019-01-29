package controllers

import (
	"github.com/astaxie/beego"
	"strings"
	"fmt"
	"gatssh/models"
)

const (
	GAT_SESSION_KEY  = "GAT_SESSION_KEY"
	GAT_SESSION_USER = "GAT_SESSION_USER"
)

var HTTPCODE = map[int]string{
	20000: "OK",
	40000: "Bad Request",
}

type baseResponse struct {
	Code int         `json:"code"`
	Msg  string      `json:"msg"`
	Data interface{} `json:"data,omitempty"`
}

type baseController struct {
	beego.Controller
	IsLogin bool
	User    string //login user
}

//Each HTTP requests going to  through here first
//Session use "user's uuid" + "request's client IP" as a way of authentication.
func (this *baseController) Prepare() {

	this.IsLogin = false

	sKey := this.GetSession(GAT_SESSION_KEY)
	sUser := this.GetSession(GAT_SESSION_USER)

	if sKey == nil || sUser == nil{
		this.routeFilter()
		return
	}
	uid := models.GetUid(sUser.(string))
	clientIp := this.getClientIp()

	if sKey.(string) != fmt.Sprint(uid+clientIp) {
		this.routeFilter()
		return
	}
	this.User = sUser.(string)
	this.IsLogin = true
	return
}

//These routers do not require session authentication.
//Because they're going to judge whether the user is logged in or not
func (this *baseController) routeFilter() {

	controllerName, _ := this.GetControllerAndAction()

	switch controllerName {
	case "UserLogin":
		return
	case "UserController":
		return
	case "GatSshQuickStart":
		return
	case "QueryGatSshTaskResults":
		return
	default:
		this.Redirect("/login", 302)
		this.ServeJSON(40000, nil)
		return
	}
	return
}

func (this *baseController) getClientIp() string {
	s := strings.Split(this.Ctx.Request.RemoteAddr, ":")
	return s[0]
}

func (this *baseController) ServeJSON(code int, data interface{}) {
	msg := baseResponse{
		Code: code,
		Msg:  HTTPCODE[code],
	}
	if data != nil {
		msg.Data = data
	}
	this.Data["json"] = msg
	this.Controller.ServeJSON()
}

