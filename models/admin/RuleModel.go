package admin

// 权限模型
import (
	"github.com/beego/beego/v2/client/orm"
	"github.com/beego/beego/v2/core/logs"
	beego "github.com/beego/beego/v2/server/web"
)

// 权限组和用户对应表
type AuthGroupAccessModel struct {
	Id      int64 `orm:"pk"`
	Uid     int64
	GroupId int
}

// 权限组表
type AuthGroupModel struct {
	Id          int `orm:"pk"`
	Title       string
	Description string
	Status      int8
	Rules       string
}

// 规则表
type RuleModel struct {
	Id         int    `orm:"pk"`
	Name       string `valid:"Required"` // 规则名
	Rule       string `valid:"Required"` // 规则，一般是路由。如： /user/list
	Param      string // 规则参数。如： limit=10&page=1
	Status     int8   `valid:"Range(0, 1)"` // 状态 1-启用 0-禁用
	CreateTime int64
	Soft       int   // 排序
	Open       int8  `valid:"Range(0, 1)"` // 是否对所有人员开放
	AddStaff   int64 // 添加人员
}

func init() {
	//设置表前缀并且注册模型
	dbPrefix, err := beego.AppConfig.String("db::dbPrefix")
	if err != nil {
		logs.Error(err)
	}
	orm.RegisterModelWithPrefix(dbPrefix, new(RuleModel), new(AuthGroupAccessModel), new(AuthGroupModel))
}

//自定义表名
func (Rule *AuthGroupAccessModel) TableName() string {
	return "auth_group_access"
}

//自定义表名
func (Rule *AuthGroupModel) TableName() string {
	return "auth_group"
}

//自定义表名
func (Rule *RuleModel) TableName() string {
	return "rule"
}
