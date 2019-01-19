package utils

import (
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/logs"
)

var GatLog = beego.BeeLogger

func initLog() {

	GatLog.Reset()
	logConfig := `{"filename":"gatssh.log","maxdays":7,"perm": "0644"}`
	if err := GatLog.SetLogger(logs.AdapterFile, logConfig); err != nil {
		panic(err)
	}
	GatLog.EnableFuncCallDepth(true)
	GatLog.SetLogger("console", "")
	GatLog.SetLevel(logs.LevelDebug)

	beego.BConfig.Log.AccessLogs = true
	beego.BConfig.Log.FileLineNum = true
}

func init()  {
	initLog()
}