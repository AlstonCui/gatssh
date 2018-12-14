package controllers

import (
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/logs"
)

var log = beego.BeeLogger

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

func init() {
	initLog()
}

func initLog() {

	log.Reset()
	logConfig := `{"filename":"gatlin.log","maxdays":7,"perm": "0644"}`
	if err := log.SetLogger(logs.AdapterFile, logConfig); err != nil {
		panic(err)
	}
	log.EnableFuncCallDepth(true)
	log.SetLogger("console", "")
	log.SetLevel(logs.LevelInformational)
	log.EnableFuncCallDepth(true)

	beego.BConfig.Log.AccessLogs = true
	beego.BConfig.Log.FileLineNum = true
}
