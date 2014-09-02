//
// Created by weibaohui on 14-9-2.
//
//文件作用:
package task
import (

	"github.com/astaxie/beego/toolbox"
	"github.com/astaxie/beego/httplib"
	"github.com/astaxie/beego"
	"fmt"
)

func init(){
	url:=beego.AppConfig.String("httpaddr")
	port:=beego.AppConfig.String("httpport")
	fetch:=fmt.Sprintf("http://%s:%s/redis/sendemail",url,port)
	tk1 := toolbox.NewTask("SendConfig", "0 */120 * * * *", func() error {httplib.Get(fetch).String(); return nil })
	toolbox.AddTask("SendConfig", tk1)

	toolbox.StartTask()
	defer toolbox.StopTask()
}
