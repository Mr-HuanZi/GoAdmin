package admin

import (
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/orm"
)

type SystemLogsModel struct {
	Id         int64 `orm:"pk"`
	Type       int8  `valid:"Range(1,2)"`
	CreateTime int64
	Ip         string
	Url        string
	UserId     int64
	Username   string
	Content    string
	Param      string
}

func init() {
	//设置表前缀并且注册模型
	orm.RegisterModelWithPrefix(beego.AppConfig.String("db::dbPrefix"), new(SystemLogsModel))
}

//自定义表名
func (SystemLogs *SystemLogsModel) TableName() string {
	return "system_logs"
}
