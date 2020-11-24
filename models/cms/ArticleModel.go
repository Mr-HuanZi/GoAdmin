package cms

import (
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/orm"
)

// 文章表模型
type ArticleModel struct {
	Id            int64  `orm:"pk"`
	Title         string `description:"标题"`
	CategoryId    int    `description:"栏目ID"`
	CategoryName  string `orm:"-"` // 栏目名，非数据库字段
	Describe      string `description:"文章描述"`
	Content       string `description:"正文"`
	Status        int8   `description:"文章状态 1-可用 2-禁用 3-删除"`
	CreateTime    int64  `description:"文章创建时间"`
	UpdateTime    int64  `description:"文章更新时间"`
	Tag           string `description:"Tags"`
	PostHits      int64  `description:"查看数"`
	PostLike      int64  `description:"点赞数"`
	CommentCount  int64  `description:"评论数"`
	CommentStatus int    `description:"评论状态;1:允许;0:不允许"`
	More          string `description:"扩展属性,如缩略图;格式为json"`
	Source        string `description:"来源"`
	Author        string `description:"作者"`
	StaffId       int64  `description:"作者ID"`
}

//自定义表名
func (Article *ArticleModel) TableName() string {
	return "article"
}

func init() {
	//设置表前缀并且注册模型
	orm.RegisterModelWithPrefix(beego.AppConfig.String("db::dbPrefix"), new(ArticleModel))
}
