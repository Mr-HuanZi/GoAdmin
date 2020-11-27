package cms

import (
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/orm"
)

// 文章表模型
type ArticleModel struct {
	Id            int64  `orm:"pk" json:"id,string"`
	Title         string `description:"标题" json:"title"`
	CategoryId    int    `description:"栏目ID" json:"category_id,string"`
	CategoryName  string `orm:"-" json:"category_name"` // 栏目名，非数据库字段
	Describe      string `description:"文章描述" json:"describe"`
	Content       string `description:"正文" json:"content"`
	Status        int8   `description:"文章状态 1-可用 2-禁用 3-删除" json:"status"`
	CreateTime    int64  `description:"文章创建时间" json:"create_time"`
	UpdateTime    int64  `description:"文章更新时间" json:"update_time"`
	Tag           string `description:"Tags" json:"tag"`
	PostHits      int64  `description:"查看数" json:"post_hits"`
	PostLike      int64  `description:"点赞数" json:"post_like"`
	CommentCount  int64  `description:"评论数" json:"comment_count"`
	CommentStatus int    `description:"评论状态;1:允许;0:不允许" json:"comment_status"`
	More          string `description:"扩展属性,如缩略图;格式为json" json:"more"`
	Source        string `description:"来源" json:"source"`
	Author        string `description:"作者" json:"author"`
	StaffId       int64  `description:"作者ID" json:"staff_id"`
	Recommend     int8   `description:"推荐位 000-默认 001-置顶 010-推荐" json:"recommend"`
}

//自定义表名
func (Article *ArticleModel) TableName() string {
	return "article"
}

func init() {
	//设置表前缀并且注册模型
	orm.RegisterModelWithPrefix(beego.AppConfig.String("db::dbPrefix"), new(ArticleModel))
}
