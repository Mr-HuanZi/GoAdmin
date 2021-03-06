package models

import (
	"github.com/beego/beego/v2/client/orm"
	"github.com/beego/beego/v2/core/logs"
	beego "github.com/beego/beego/v2/server/web"
)

type AssetModel struct {
	Id         int64
	Path       string
	Md5        string
	Sha1       string
	Status     int8
	FileName   string
	FileInfo   string
	CreateTime int64
	AddStaff   int64
	Thumb      string
	Size       int64
}

func init() {
	//设置表前缀并且注册模型
	dbPrefix, err := beego.AppConfig.String("db::dbPrefix")
	if err != nil {
		logs.Error(err)
	}
	orm.RegisterModelWithPrefix(dbPrefix, new(AssetModel))
}

//自定义表名
func (model *AssetModel) TableName() string {
	return "asset"
}
