package routers

import (
	"gatlin/controllers"
	"github.com/astaxie/beego"
)

func init() {
	beego.Router("/v1/remoteCmd", &controllers.GatSshMultiShoot{})
	//beego.Router("/login", &controllers.UserController{}, `get:PageLogin`)
	//beego.Router("/register", &controllers.UserController{}, `post:Register`)
	//beego.Router("/reallogin", &controllers.UserController{}, `post:Reallogin`)
	//beego.Router("/v1/remoteShell",&controllers.RemoteShell{})

	beego.Router("/",&controllers.MainController{})
	beego.Router("/login",&controllers.UserLogin{})
	beego.AutoRouter(&controllers.UserController{})


}
