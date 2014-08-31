package controllers

import (

	"github.com/astaxie/beego"
	. "UnionGo/models"
	"github.com/astaxie/beego/orm"
	. "UnionGo/Library"

)

type OptionController struct {
	BaseController

}

func (this *OptionController) Index() {

	ss := GetOptions()
	beego.Debug(ss["key"])
	this.TplNames = "option.html"
	this.Render()
}

func (this *OptionController) Save() {
 	data := `{"list":` + this.GetString("data") + `}`
	h := new(Option)
	diy:=this.GetUserInfo()
	h.SaveList(data,diy)

	this.Data["json"] = "ok"
	this.ServeJson()
}



func (this *OptionController) Get() {



	//
//	pageIndex	0
//	pageSize	10
//	sortField
//	sortOrder
	this.GetString("pageIndex")
	var pulist []Option
	o := orm.NewOrm()
	pu := new(Option)
	qs := o.QueryTable(pu)
	qs = qs.Limit(20, 0)
	qs.All(&pulist)
	this.Data["json"] = &MiniuiGrid{1000, &pulist}
	this.ServeJson()

}
