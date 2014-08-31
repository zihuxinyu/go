package main

import (
	"github.com/astaxie/beego"
	_ "ShadowManage/routers"
	//. "github.com/zihuxinyu/GoLibrary"
)

func main() {
	beego.SessionOn = true


//	white := []string{
//		"/login",
//		"/redis/sendemail",
//		"/redis/index",
//		"/login/index",
//	}


//	var FilterUser = func(ctx *context.Context) {
//		beego.Debug("ddd", ctx.Request.RequestURI)
//		value := ctx.Input.Session("User")
//		if value != "Manager" && !StringsContains(white, ctx.Request.RequestURI) {
//			ctx.Redirect(302, "/login/index")
//		}
//	}
//	beego.InsertFilter("/redis/*", beego.BeforeExec, FilterUser)

	beego.Run()
}

