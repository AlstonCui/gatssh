package routers

import (
	"gatssh/controllers"
	"github.com/astaxie/beego"
)

func init() {

	beego.Router("/",&controllers.MainController{})

	beego.Router("/login",&controllers.UserLogin{},"Get:Login")
	beego.Router("/loginAuth",&controllers.UserLogin{},"Post:LoginAuth")

	beego.Router("/v1/quickStart",&controllers.GatSshQuickStart{},"post:QuickStart")
	beego.Router("/v1/StartReceiveFormWS",&controllers.GatSshQuickStart{},"get:StartSendByWS")
	beego.Router("/v1/queryTaskDetails",&controllers.QueryGatSshTaskDetails{})

}
