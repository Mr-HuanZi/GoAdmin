package admin

import (
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/orm"
)

type AuthGroupAccessModel struct {
	Id      int64 `orm:"pk"`
	Uid     int64
	GroupId int
}

func init() {
	//设置表前缀并且注册模型
	orm.RegisterModelWithPrefix(beego.AppConfig.String("db::dbPrefix"), new(AuthGroupAccessModel))
}

//自定义表名
func (Rule *AuthGroupAccessModel) TableName() string {
	return "auth_group_access"
}
