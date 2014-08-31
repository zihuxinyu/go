package portal_user

import (
	"github.com/astaxie/beego/orm"
	"time"
	. "UnionGo/Library"
)

type Portal_user struct{
	Guid                   int `orm:"pk" `
	User_code              string
	User_name              string
	Dpt_name               string
	Msgexpdate             time.Time `orm:"auto_now_add;type(datetime)"`

}

func init() {
	orm.RegisterModel(new(Portal_user))

	ModelCache.Set("Portal_user", func() interface{} {return &Portal_user{}})

}
func (m *Portal_user) TableName() string {
	return TableName("portal_user")
}
func (h Portal_user) SaveList(data string, diy interface{}) {

	SaveMiniUIData("Portal_user", data, diy)
}
