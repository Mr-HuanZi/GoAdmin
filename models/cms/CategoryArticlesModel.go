package cms

import (
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/orm"
)

type CategoryArticlesModel struct {
	Id int64 `orm:"pk"`
	ArticlesId int64
	CategoryId int
	ListOrder float64
	Status int8
}

//自定义表名
func (Category *CategoryArticlesModel) TableName() string {
	return "category_articles"
}

func init() {
	//设置表前缀并且注册模型
	orm.RegisterModelWithPrefix(beego.AppConfig.String("db::dbPrefix"), new(CategoryArticlesModel))
}

func UpdateCategoryArticles(article ArticleModel, category CategoryModel) error {
	o := orm.NewOrm()
	// 删除旧数据
	_, delErr := o.Delete(&CategoryArticlesModel{ArticlesId:article.Id})
	if delErr != nil {
		return delErr
	}
	// 插入新数据
	_, insErr := o.Insert(&CategoryArticlesModel{
		ArticlesId: article.Id,
		CategoryId: category.Id,
		ListOrder:  0,
		Status:     article.Status,
	})
	if insErr != nil {
		return insErr
	}
	return nil
}