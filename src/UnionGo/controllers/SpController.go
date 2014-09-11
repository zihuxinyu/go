//异常码(AABB)	含义
//100	参数错误
//101	方法参数缺失
//102	方法参数错误
//103	鉴权参数缺失
//104	鉴权参数错误
//105	sessionkey参数缺失
//106	sessionkey参数错误
//107	appkey参数缺失
//108	appkey参数错误
//109	时间戳参数缺失
//110	时间戳参数错误
//111	版本参数缺失
//112	版本参数错误
//113	文件未找到
//114	文件类型错误
//115	HTTP请求错误
//116	API系统错误
//117	API系统参数错误
//118	非法访问
//119	权限错误
//120	返回传输错误
//121	数字转换错误
//122	数据库访问错误
//123	无权限
//124	数字转换错误
//125	参数少于规定
//126	参数多余规定
//10010003	用户未登录
//10010004	用户无访问权限
//10010005	该地市无产品

package controllers

import (
	. "UnionGo/Library"
	"crypto/sha1"
	"fmt"
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/httplib"
	"sort"
	. "UnionGo/models/weixin"
	"UnionGo/models/weixin"
)

var (
	appkey      = "dy-weixinyx"
	securityKey = "dy-weixinyx"
	loginName   = "dy-weixinyx"
	password    = "123456"
	version     = "1.0"
	serverurl   = "http://112.231.23.40:8080/API/api"
)

type SpController struct {
	beego.Controller
	BaseController
}

func (this *SpController) Index() {

	DoTest()
	this.TplNames = "index.html"
	this.Data["string"] = "ddd"
	this.RenderString()

}

func DoTest() {

	ss,err := productService_categories()
	fmt.Println(ss,err)

}

//产品分类
//名称	类型	是否必须	描述
//loginName	String	N	代理商登录名
//parent	String	N	产品类别ID

func productService_categories() (content string, err error) {
	method := "productService.categories"
	Params := map[string]interface{}{
		"method": method,
	}
	content,err =Post2Serv(Params)
	m:=Splog{
		Method:method,
		Params:fmt.Sprintf("%s", Params),
		Content:fmt.Sprintf("%s", content),
		Err:fmt.Sprintf("%s", err),
	}
	weixin.AddSplog(&m)
	return
}

//佣金列表
//名称	类型	是否必须	描述
//loginName	String	Y	代理商登录名
//billMonth	Sring	N	佣金月份(格式：YYYYMM)
//start	Int	N	开始行
//limit	Int	N	行数
//todo:目前少参数，要新版API
func commissionService_list(start, limit string) (content string, err error) {
	method := "commissionService.list"
	Params := map[string]interface{}{
		"method":    method,
		"billMonth": "201408",
		"start":     start,
		"limit":     limit,
	}

	return Post2Serv(Params)
}

//已推荐产品产品列表详情，针对代理商
//名称	类型	是否必须	描述
//loginName	String	Y	代理商登录名
//startTime	Datetime	N	开始日期，注：大于或等于开始时间 格式yyyyMMdd
//endTime	Datetime	N	结束日期，注：小于结束时间 yyyyMMdd
//status	String	N	参数值：成功状态(success);失败状态(fail);正在处理状态(doing);不传表示查询所有状态的数据.
//orderId	String	N	请求推荐之后,api返回的推荐流水号(推荐ID)
//start	Int	Y	开始行
//limit	Int	Y	行数
func recommendService_detail(start, limit string) (content string, err error) {
	method := "recommendService.detail"
	Params := map[string]interface{}{
		"method": method,
		"start":  start,
		"limit":  limit,
	}
	return Post2Serv(Params)
}

//推荐
//名称	类型	是否必须	描述
//loginName	String	Y	代理商登录名
//mobile	String	Y	被推荐人手机号码
//productId	String	Y	产品ID
//effectiveMode	String	N	生效方式：0 当月开通(默认) 1 次月开通
func recommendService_recommend(mobile, productId, effectiveMode string) (content string, err error) {
	method := "recommendService.recommend"
	Params := map[string]interface{}{
		"method":        method,
		"mobile":        mobile,
		"productId":     productId,
		"effectiveMode": effectiveMode,
	}
	return Post2Serv(Params)
}

//产品列表
//名称	类型	是否必须	描述
//loginName	String	Y	代理商登录名
//key	String	N	产品名称关键字
//productCode	String	N	产品ID
//category	String	N	产品分类ID
//start	int	Y	分页起始页
//limit	int	Y	每页记录条数
func productService_list(start, limit string) (content string, err error) {
	method := "productService.list"
	Params := map[string]interface{}{

		"method": method,
		"start":  start,
		"limit":  limit,
	}
	return Post2Serv(Params)
}

//单个产品查询
func productService_Single(productCode, start, limit string) (content string, err error) {
	method := "productService.list"
	Params := map[string]interface{}{

		"method":      method,
		"start":       start,
		"limit":       limit,
		"productCode": productCode,
	}
	return Post2Serv(Params)
}

//已订购列表，针对被推荐的用户
//名称	类型	是否必须	描述
//mobile	String	Y	查询用户的联通手机号
func userService_getOrders(mobile string) (content string, err error) {
	method := "userService.getOrders"
	Params := map[string]interface{}{

		"method": method,
		"mobile": mobile,
	}
	return Post2Serv(Params)
}

//获取用户信息
//名称	类型	是否必须	描述
//mobile	String	Y	用户手机号

func userService_userInfo(mobile string) (content string, err error) {
	Params := map[string]interface{}{
		"mobile": mobile,
		"method": "userService.userInfo",
	}
	return Post2Serv(Params)
}

//执行API
func Post2Serv(Params map[string]interface{}) (content string, err error) {
	Params["appkey"] = appkey
	Params["v"] = version
	Params["loginName"] = loginName
	sig := Signature(Params, securityKey)

	b := httplib.Post(serverurl)
	for k, v := range Params {
		b.Param(k, v.(string))
	}
	//加上签名
	b.Param("sig", sig)

	str, err := b.String()
	if err != nil {
		beego.Debug(err)
	}
	return str, err
}

//签名
//所有参数key按升序排列，顺序取出按key+value组合在一起，最后加上securitykey，做sha1运算
//返回运算值作为sig
func Signature(Params map[string]interface{}, securityKey string) string {

	//	Params := map[string]interface {}{
	//
	//		"uName":"13521081739",
	//		"uPass":"123456",
	//		"method":"user.login",
	//		"format":"json",
	//		"version":"1.0",
	//		"appkey":"a6479ba4c45b658c",
	//	}

	//appkey:="dy-weixinyx"
	//securityKey := "fksds2323dsdf"

	strs := sort.StringSlice(MapKeys(Params))
	sort.Strings(strs)
	str := ""
	for _, s := range strs {
		str += s + Params[s].(string)
	}
	str += securityKey
	//fmt.Printf("Signature()\n", str)
	h := sha1.New()
	h.Write([]byte(str))
	return fmt.Sprintf("%x", h.Sum(nil))
}
