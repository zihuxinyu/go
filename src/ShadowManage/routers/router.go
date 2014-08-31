package routers

import (
	"github.com/astaxie/beego"
	"ShadowManage/controllers"
)

func init() {
	beego.AutoRouter(&controllers.RedisController{})
	beego.AutoRouter(&controllers.LoginController{})
}
