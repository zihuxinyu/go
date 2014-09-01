package controllers

import (
	. "github.com/zihuxinyu/GoLibrary"
	"fmt"
	"encoding/json"
	"time"
	"reflect"
	"strings"
	. "github.com/astaxie/beego/logs"
	"github.com/astaxie/beego"
)

type RedisController struct {
	beego.Controller
	BaseController

}

func Notify() {
	go ExecSystem("ps -ef|grep ./server|grep -v grep|awk '{printf $2}'|xargs kill -1")
}

func (this *RedisController) SendEmail() {
	log := NewLogger(10000)
	log.SetLogger("smtp", `{"username":"aixinit@126.com","password":"9Loveme?","host":"smtp.126.com:25","sendTos":["491923016@qq.com"]}`)
	StoragePtr = NewStorage()

	Userlist, _ := StoragePtr.GetList()
	data, _ := json.Marshal(Userlist)
	log.Critical("back", string(data), this.getClientIp())

}

func (this *RedisController) Index() {



	this.TplNames = "redis.html"
	this.Render()
}

func (this *RedisController) Save() {

	StoragePtr = NewStorage()

	data := `{"list":` + this.GetString("data") + `}`

	//整理为可识别格式

	var usr User
	var dataList DataList
	json.Unmarshal([]byte(data), &dataList)

	StructType := reflect.TypeOf(usr) //通过反射获取type定义


	//按struct 遍历得到定义，及得到的值
	for _, SingleItem := range dataList.List {
		fmt.Println("第一次", SingleItem)
		for i := 0; i < StructType.NumField(); i++ {
			f := StructType.Field(i)

			if SingleItem[f.Name] != nil {
				//								if f.Name == "Port" {
				//									switch SingleItem[f.Name].(type){
				//
				//									case float64:
				//										fmt.Println("float",SingleItem[f.Name])
				//										limit, _ := strconv.Itoa(SingleItem[f.Name].(int))
				//
				//										SingleItem[f.Name] = float64(limit)
				//									}
				//
				//								}
				//fmt.Println("格式整理",f.Name,reflect.TypeOf( SingleItem[f.Name]))
				if f.Type == reflect.TypeOf(time.Now()) {
					//对时间格式进行特殊的处理，进行时区转换，miniui过来的时间加+08:00
					//处理为go转换string为时间需要的标准时间格式
					ss := fmt.Sprintf("%s", SingleItem[f.Name])
					ss = strings.Replace(ss, "T", " ", -1)
					if !strings.Contains(ss, "+08:00") {
						ss = ss+" +08:00"
					}
					//fmt.Println("时间格式整理"+ss)
					ss = strings.Replace(ss, "+08:00", " +08:00", -1)

					t, _ := time.Parse("2006-01-02 15:04:05 -07:00 ", ss)


					//转换正确的时间回填
					//SingleItem[f.Name] = t.Format("2006-01-02 15:04:05")
					SingleItem[f.Name] = t
				}

			}
			fmt.Println(f.Name, SingleItem[f.Name], reflect.TypeOf(SingleItem[f.Name]))
		}
		fmt.Println("第二次", SingleItem)

		if state := SingleItem["_state"]; state != "" {
			//先将map 对应为json
			x, _ := json.Marshal(SingleItem)
			fmt.Println(reflect.TypeOf(x), string(x))
			//再将json对应为struct
			json.Unmarshal(x, &usr)

			switch state {
			case "modified":
				fmt.Println("modified", &usr)
				if err := StoragePtr.Set("user:"+usr.Name, usr); err != nil {
					fmt.Println("修改错误", err)
				}
			case "added":
				fmt.Println("added", &usr)
				if err := StoragePtr.Set("user:"+usr.Name, usr); err != nil {
					fmt.Println("新增错误", err)
				}
			case "removed":
				fmt.Println("removed", &usr)
				if err := StoragePtr.Del("user:" + usr.Name); err != nil {
					fmt.Println("删除错误", err)
				}
			}
		}
	}
	Notify()//修改后通知服务器进行更新
	this.Data["json"] = "ok"
	this.ServeJson()

}


func (this *RedisController) Get() {

	StoragePtr = NewStorage()

	Userlist, _ := StoragePtr.GetList()


	this.Data["json"] = &MiniuiGrid{int64(len(Userlist)), &Userlist}
	this.ServeJson()



}

func (this *RedisController) Reset(){
	StoragePtr = NewStorage()
	Username:=this.GetString("user")
	if err:=StoragePtr.ResetUsed("flow:"+Username);err!=nil{
		this.Data["json"]=err
	}else{
		this.Data["json"]="已完成"
	}

	this.ServeJson()
}
//获取用户的限制信息，包括当前用量，限制记录
func (this *RedisController) GetLimit() {

	StoragePtr = NewStorage()

	Username:=this.GetString("user")

	size,_:=StoragePtr.GetSize("flow:"+Username)
	log,_:=StoragePtr.GetLog("log:"+Username)
	type Uinfo struct {
		Name string
		Size int64
		Log string
	}
	uinfo:=Uinfo{Name:Username,Size:size,Log:log}
	this.Data["json"] = &uinfo
	this.ServeJson()

}
