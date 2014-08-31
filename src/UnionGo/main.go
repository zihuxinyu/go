package main

import (
	"UnionGo/controllers"
	_ "UnionGo/routers"
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/orm"
	_ "github.com/go-sql-driver/mysql" // import your used driver
	"fmt"
	//"github.com/astaxie/beego/context"
)

func init() {


	dbPass, _ := beego.GetConfig("string", "db_pass")
	dbUser, _ := beego.GetConfig("string", "db_user")
	dbHost, _ := beego.GetConfig("string", "db_host")
	dbPort, _ := beego.GetConfig("string", "db_port")
	dbName, _ := beego.GetConfig("string", "db_name")

	maxIdleConn, _ := beego.GetConfig("int", "db_max_idle_conn")
	maxOpenConn, _ := beego.GetConfig("int", "db_max_open_conn")
	dbLink := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8&loc=%s", dbUser, dbPass, dbHost, dbPort, dbName,"Asia%2FChongqing")
	orm.RegisterDriver("mysql", orm.DR_MySQL)
	orm.RegisterDataBase("default", "mysql", dbLink, maxIdleConn.(int), maxOpenConn.(int))
}

func main() {
	orm.Debug = true
	beego.SessionOn = true

//	var FilterUser = func(ctx *context.Context) {
//		beego.Debug("ddd",ctx.Request.RequestURI)
//		_, ok := ctx.Input.Session("uid").(int)
//		if !ok && ctx.Request.RequestURI != "/login" {
//			ctx.Redirect(302, "/login")
//		}
//	}
//
//	beego.InsertFilter("/d/*",beego.BeforeExec,FilterUser)
	beego.AutoRouter(&controllers.MainController{})
	beego.AutoRouter(&controllers.OptionController{})
	beego.Run()

}
