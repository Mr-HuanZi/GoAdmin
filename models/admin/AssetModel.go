package admin

import (
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/orm"
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
	orm.RegisterModelWithPrefix(beego.AppConfig.String("db::dbPrefix"), new(AssetModel))
}

//自定义表名
func (model *AssetModel) TableName() string {
	return "asset"
}
