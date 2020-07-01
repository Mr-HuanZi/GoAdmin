package admin

import (
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/orm"
)

type RuleModel struct {
	Id         int    `orm:"pk"`
	Name       string // 规则名
	Rule       string // 规则，一般是路由。如： /user/list
	Param      string // 规则参数。如： limit=10&page=1
	Status     uint8  // 状态 1-启用 0-禁用
	CreateTime int64
	Soft       int   // 排序
	Open       int8  // 是否对所有人员开放
	AddStaff   int64 // 添加人员
}

func init() {
	//设置表前缀并且注册模型
	orm.RegisterModelWithPrefix(beego.AppConfig.String("db::dbPrefix"), new(RuleModel))
}

//自定义表名
func (User *RuleModel) TableName() string {
	return "rule"
}
