package Library

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/astaxie/beego/orm"
	"reflect"
	"strings"
	"time"
	"github.com/astaxie/beego"
)

//返回带前缀的表名
func TableName(str string) string {
	return fmt.Sprintf("%s%s", beego.AppConfig.String("dbprefix"), str)
}

type MiniuiGrid struct {
	Total int64 `json:"total"`
	Data  interface{} `json:"data"`
}

//构造新的struct,接收json 绑定
//	[{"_state":"modified","Guid":22,"Msgexpdate":"2014-08-15T15:56:18"},{"_state":"modified","Guid":23,"Msgexpdate":"2014-08-15T15:56:10"}]

type DataList struct {
	List [] map[string]interface{}
}

//找出beegoModel的主键
func GetModelPk(obj interface{}) (pkFiledName string) {
	s := reflect.TypeOf(obj).Elem() //通过反射获取type定义
	for i := 0; i < s.NumField(); i++ {
		pkFiled := s.Field(i)
		tags := strings.Split(pkFiled.Tag.Get("orm"), ",")
		//fmt.Println(tags)
		for _, v := range tags {
			if strings.ContainsAny(v, ";") {
				for _, vv := range strings.Split(v, ";") {

					if strings.ToLower(vv) == "pk" {
						pkFiledName = pkFiled.Name
						break
					}
				}
			}else {
				if strings.ToLower(v) == "pk" {
					pkFiledName = pkFiled.Name
					break
				}
			}
		}
		//得到值就退出循环
		if len(pkFiledName) > 0 {
			break
		}
	}
	return pkFiledName
}

//格式化miniui过来的时间,得到修改过的字段放到m中
func MiniUIDataUpdate(obj interface{}, SingleItem map[string]interface{}, m orm.Params, state string , diys interface{}) error {


	if !isPtr(obj) {
		return errors.New(fmt.Sprintf("只支持指针类型，不支持`%T`", obj))
	}
	StructType := reflect.TypeOf(obj).Elem() //通过反射获取type定义

	//进行type转换
	var diy map[string]interface{}
	diy = diys.(map[string]interface{})



	if state == "modified" {
		//		//修改时删除创建人得信息，因为修改时不能改变创建人信息
		delete(diy, "Creatorid")
		delete(diy, "Createdate")
	}

	for i := 0; i < StructType.NumField(); i++ {
		f := StructType.Field(i)
		//fmt.Println(f.Name, f.Type, reflect.TypeOf( v[f.Name]))

		//此处日后可以做根据字段名注入逻辑控制后的字段，比如创建人等信息
		if SingleItem[f.Name] != nil {
			//fmt.Println("格式整理",f.Name,reflect.TypeOf( SingleItem[f.Name]))
			if f.Type == reflect.TypeOf(time.Now()) {
				//对时间格式进行特殊的处理，进行时区转换，miniui过来的时间加+08:00
				//处理为go转换string为时间需要的标准时间格式
				ss := fmt.Sprintf("%s", SingleItem[f.Name])
				ss = strings.Replace(ss, "T", " ", -1)

				if !strings.Contains(ss,"+08:00") {
					ss = ss+" +08:00"
				}
				//fmt.Println("时间格式整理"+ss)
				ss = strings.Replace(ss, "+08:00", " +08:00", -1)

				t, _ := time.Parse("2006-01-02 15:04:05 -07:00 ", ss)

				m[f.Name] = t
				//转换正确的时间回填
				//SingleItem[f.Name] = t.Format("2006-01-02 15:04:05")
				SingleItem[f.Name] = t
			} else {

				m[f.Name] = SingleItem[f.Name]
			}
			//fmt.Println("格式整理后",f.Name,reflect.TypeOf( SingleItem[f.Name]))

		}

		for x := range diy {
			//fmt.Println(x,diy[x],f.Name,reflect.TypeOf(SingleItem[f.Name]))
			if x == f.Name {
				m[f.Name] = diy[x]
				SingleItem[f.Name] = diy[x]
			}
		}

	}

	//先将map 对应为json
	x, _ := json.Marshal(SingleItem)
	//fmt.Println(reflect.TypeOf(x), string(x))
	//再将json对应为struct
	json.Unmarshal(x, obj)

	return nil
}

///将miniui过来的数据保存，根据ModelName通过ModelCache模块获得model实例
//diy:用户信息，包括修改人修改时间等
func SaveMiniUIData(ModelName string, data string, diy interface{}) error {
	//根据名称获取Model
	reflecty, _ := ModelCache.Get(ModelName)
	//得到实例
	reflectx := reflecty()



	//整理为可识别格式
	var dataList DataList
	json.Unmarshal([]byte(data), &dataList)


	//按struct 遍历得到定义，及得到的值
	for _, SingleItem := range dataList.List {

		if state := SingleItem["_state"]; state != "" {

			Params := make(orm.Params)
			//格式化miniui过来的时间,得到修改过的字段放到m中
			MiniUIDataUpdate(reflectx, SingleItem, Params, state.(string), diy)

			switch state {
			case "modified":
				pk := GetModelPk(reflectx)
				//获取表名，要求必须有TableName方法
				if tablenameMC := reflect.ValueOf(reflectx).MethodByName("TableName"); tablenameMC.IsValid() {
					tablename := tablenameMC.Call(nil)[0].Interface().(string)
					//fmt.Println("dddd",tablename.Interface().(string))
					//orm.NewOrm().QueryTable(reflect.TypeOf(reflectx).Elem().Name()).Filter(pk, SingleItem[pk]).Update(Params)
					orm.NewOrm().QueryTable(tablename).Filter(pk, SingleItem[pk]).Update(Params)

				}else {
					return errors.New(fmt.Sprintf("获取表名，要求必须有TableName方法"))

				}

			case "added":
				orm.NewOrm().Insert(reflectx)
			case "removed":

				orm.NewOrm().Delete(reflectx)

			}
		}
	}
	return nil
}
