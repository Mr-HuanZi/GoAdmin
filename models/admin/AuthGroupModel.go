package admin

import (
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/orm"
)

type AuthGroupModel struct {
	Id          int `orm:"pk"`
	Title       string
	Description string
	Status      int8
	Rules       string
}

func init() {
	//设置表前缀并且注册模型
	orm.RegisterModelWithPrefix(beego.AppConfig.String("db::dbPrefix"), new(AuthGroupModel))
}

//自定义表名
func (Rule *AuthGroupModel) TableName() string {
	return "auth_group"
}
