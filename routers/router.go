package routers

import (
	"gatlin/controllers"
	"github.com/astaxie/beego"
)

func init() {
	beego.Router("/v1/remoteCmd", &controllers.GatSshOneShoot{})
	//beego.Router("/v1/remoteShell",&controllers.RemoteShell{})
}
