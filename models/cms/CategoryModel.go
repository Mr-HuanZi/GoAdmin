package cms

import (
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/logs"
	"github.com/astaxie/beego/orm"
)

//栏目表模型
type CategoryModel struct {
	Id int `orm:"pk"`
	ParentId int `description:"栏目父id"`
	Status int8 `description:"状态;1:显示;2:隐藏" valid:"Range(0,2)"`
	Sort float32 `description:"排序"`
	Name string `description:"分类名称" valid:"Required"`
	Description string `description:"分类描述"` 
	Alias string `description:"分类别名" valid:"Required;AlphaNumeric"`
	ListTpl string `description:"分类列表模板"` 
	OneTpl string `description:"分类文章页模板"` 
	CreateTime int64 `description:"创建时间"` 
	Icon string `description:"分类图标"` 
	Thumbnail int `description:"分类封面图"` 
	More string `description:"扩展属性,格式为json"` 
}

//自定义表名
func (Category *CategoryModel) TableName() string {
	return "category"
}

func init() {
	//设置表前缀并且注册模型
	orm.RegisterModelWithPrefix(beego.AppConfig.String("db::dbPrefix"), new(CategoryModel))
}

func UpdateCategory(data *CategoryModel) (int64, error) {
	o := orm.NewOrm()
	//查找对应ID的数据是否存在
	cate := CategoryModel{Id:data.Id}
	err := o.Read(&cate)
	if err == orm.ErrNoRows {
		// 没有相应记录
		return 0, err
	} else if err == orm.ErrMissPK {
		// 主键丢失
		logs.Error(orm.ErrMissPK)
		return 0, err
	} else {
		logs.Info(cate)
		if data.Status == 0 {
			data.Status = 1
		}
		cate.Alias = data.Alias
		cate.ParentId = data.ParentId
		cate.Status = data.Status
		cate.Sort = data.Sort
		cate.Name = data.Name
		cate.Description = data.Description
		cate.ListTpl = data.ListTpl
		cate.OneTpl = data.OneTpl
		cate.Icon = data.Icon
		cate.Thumbnail = data.Thumbnail
		cate.More = data.More

		if num, err := o.Update(&cate); err == nil {
			return num, nil
		} else {
			logs.Error(err)
			return 0, err
		}
	}
}