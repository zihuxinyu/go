//
// Created by weibaohui on 14-9-18.
//
//文件作用:
package controllers
import (
	"github.com/ascoders/alipay"
	"github.com/astaxie/beego"
	"html/template"
	"fmt"
)

func init() {
	//初始化支付宝插件
	alipay.AlipayPartner = "2088302274890164"
	alipay.AlipayKey = "fgaasbl27grwg53ral6oopq9qo3hulln"
	alipay.WebReturnUrl = "http://127.0.0.1:8080/alipay/AlipayReturn" //替换成你的 异步回调页面
	alipay.WebNotifyUrl = "http://127.0.0.1:8080/alipay/AlipayNotify" //替换成你的 同步回调页面
	alipay.WebSellerEmail = "f618_com@yeah.net"        //替换成你的 支付宝账号邮箱
	beego.Debug("init")
}

type AlipayController struct {
	beego.Controller

}
func (this *AlipayController) Index() {

	//创建订单order，生成了各种信息包括订单的唯一id
	//获取支付宝即时到帐的自动提交表单
	//四个参数分别是 订单唯一id(string) 充值金额(int) 实际充值的游戏币(int) 充值时给用户的描述(string)
	form := alipay.CreateAlipaySign("1234567890", 12, 12, "wodemumu-buyer", "我酷游戏-充值代金券")
	//为了更好的用户体验，可以以json方式调用，返回了json类型字符串
	this.Data["data"] = template.HTML(form)
	beego.Debug(form)
	this.TplNames="alipay.html"
	this.Render()
	//前台接收到字符串后直接输出即可跳转
	//document.write(data);

}

/* 接收支付宝同步跳转的页面 */
func (this *AlipayController) AlipayReturn() {
	//错误代码(1为成功) 订单id(使用它查询订单) 买家支付宝账号(这个不错) 支付宝id(支付宝账单id)
	status, orderId, buyerEmail, tradeNo := alipay.AlipayReturn(&this.Controller)

	fmt.Println(status,orderId, buyerEmail, tradeNo)

}

/* 被动接收支付宝异步通知的页面 */
func (this *AlipayController) AlipayNotify() {
	status, orderId, buyerEmail, tradeNo := alipay.AlipayNotify(&this.Controller)
	fmt.Println(status,orderId, buyerEmail, tradeNo)
}
