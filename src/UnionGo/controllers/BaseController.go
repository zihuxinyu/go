package controllers

import (
	"github.com/astaxie/beego"
	"encoding/json"
	"strings"
	"github.com/astaxie/beego/orm"
	. "UnionGo/models"
	. "UnionGo/Library"
)

type BaseController struct {
	beego.Controller
}

 
func (c *BaseController) GetUserId() string {
	if userID := c.GetSession("UserID"); userID != nil {
		return userID.(string)
	}
	return ""
}

// 是否已登录
func (c *BaseController) HasLogined() bool {
	return c.GetUserId() != ""
}

//是否post提交
func (this *BaseController) isPost() bool {
	return this.Ctx.Request.Method == "POST"
}
//获取用户IP地址
func (this *BaseController) getClientIp() string {
	s := strings.Split(this.Ctx.Request.RemoteAddr, ":")
	return s[0]
}

func (c *BaseController) ClearSession() {
	c.DelSession("UserID")
}

// 修改session
func (c *BaseController) UpdateSession(key, value string) {
	c.SetSession(key, value)
}

// 返回json
func (c *BaseController) Json(i interface{}) string {
	// b, _ := json.MarshalIndent(i, "", " ")
	b, _ := json.Marshal(i)
	return string(b)
}

func (c *BaseController) GetUserInfo() map[string]interface{}{
	diy := map[string]interface {}{
		//"User_name":"都是我",
//		"Creatorid":c.GetSession("User_name").(string),
//		"Createdate":TimeNowString(),
//		"Modifierid":c.GetSession("User_name").(string),
//		"Modifydate":TimeNowString(),
		"Creatorid":"weibh",
		"Createdate":TimeLocal(),

		"Modifierid":"weibh",
		"Modifydate":TimeLocal(),
	}
	return diy
}


func GetOptions() map[string]string {
	if !Cache.IsExist("options") {
		var result []*Option
		o := orm.NewOrm()
		o.QueryTable(&Option{}).All(&result)
		options := make(map[string]string)
		for _, v := range result {
			options[v.Name] = v.Value
		}
		Cache.Put("options", options, 0)
	}
	v := Cache.Get("options")
	return v.(map[string]string)
}
