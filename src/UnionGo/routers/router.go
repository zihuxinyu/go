package routers

import (
	"github.com/astaxie/beego"
	"UnionGo/controllers"
)

func init() {
//	beego.Router("/", &controllers.MainController{})


	beego.Router("/", &controllers.OptionController{})
}
