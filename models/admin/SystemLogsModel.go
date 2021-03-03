package admin

import (
	"github.com/beego/beego/v2/client/orm"
	"github.com/beego/beego/v2/core/logs"
	beego "github.com/beego/beego/v2/server/web"
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
	dbPrefix, err := beego.AppConfig.String("db::dbPrefix")
	if err != nil {
		logs.Error(err)
	}
	orm.RegisterModelWithPrefix(dbPrefix, new(SystemLogsModel))
}

//自定义表名
func (SystemLogs *SystemLogsModel) TableName() string {
	return "system_logs"
}
