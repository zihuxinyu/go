//
// Created by weibaohui on 14-8-28.
//
//文件作用:登录
package controllers

import "github.com/astaxie/beego"

type Manager struct {
	ManagerName string `form:"username"`
	ManagerPsw  string `form:"pwd"`
}

type LoginController struct {
	beego.Controller
	BaseController
}

func (this *LoginController) Index() {
	if this.Ctx.Request.Method == "POST"{
		u := Manager{}
		if err := this.ParseForm(&u); err != nil {
			beego.Debug(err)
		}
		if u.ManagerName=="lovemumu" && u.ManagerPsw=="loveme"{
			this.SetSession("User","Manager")//login
			this.Data["json"]="ok"
		}else{
			this.Data["json"]="no"
		}

		this.ServeJson()
	}

	this.TplNames = "login.html"
	this.Render()
}
