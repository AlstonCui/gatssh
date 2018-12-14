package routers

import (
	"gatlin/controllers"
	"github.com/astaxie/beego"
)

func init() {
    beego.Router("/v1/remoteCmd", &controllers.RemoteCmd{})
}
