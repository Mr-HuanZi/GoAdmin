package cms

import (
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/logs"
	"github.com/astaxie/beego/orm"
	"github.com/astaxie/beego/validation"
	"go-admin/controllers/admin"
	"go-admin/models/cms"
	"strconv"
	"time"
)

type ArticleController struct {
	admin.BaseController
}

type ArticleListSearch struct {
	Limit           int    `valid:"Range(0, 1000)"` //分页每页显示的条数
	Page            int    `valid:"Min(1)"`         //当前页码
	Title           string //文章标题
	CreateTimeStart int    //文章创建(发表)时间
	CreateTimeEnd   int    //文章创建(发表)时间
	Status          uint8  //文章状态
	Category        int    //文章分类
}

type ArticleListResult struct {
	Total int64
	List  []*cms.ArticleModel
}

func (c *ArticleController) List() {
	var (
		ArticleListSearchS ArticleListSearch
		Err                error
		Article            = new(cms.ArticleModel)
		Data               = new(ArticleListResult)
	)
	_ = c.GetRequestJson(&ArticleListSearchS, false)
	logs.Info(ArticleListSearchS)

	//获取每页记录条数
	if ArticleListSearchS.Limit <= 0 {
		limit := beego.AppConfig.String("cms::limit")
		ArticleListSearchS.Limit, Err = strconv.Atoi(limit)
		if Err != nil {
			logs.Error(Err)
			c.Response(500, "", nil)
		}
	}

	//页码
	if ArticleListSearchS.Page <= 0 {
		ArticleListSearchS.Page = 1
	}

	//计算页码偏移量
	offset := (ArticleListSearchS.Page - 1) * ArticleListSearchS.Limit
	o := orm.NewOrm()
	qs := o.QueryTable(Article)

	//状态搜索
	if ArticleListSearchS.Status != 0 {
		qs.Filter("status", ArticleListSearchS.Status)
	} else {
		qs.Filter("status__in", 1, 2, 3)
	}

	//开始时间
	if ArticleListSearchS.CreateTimeStart != 0 {
		qs.Filter("create_time__gte", ArticleListSearchS.CreateTimeStart)
	}

	//结束时间
	if ArticleListSearchS.CreateTimeEnd != 0 {
		qs.Filter("create_time__lte", ArticleListSearchS.CreateTimeEnd)
	}

	// 获取总条目
	cnt, errCount := qs.Count()
	if errCount != nil {
		logs.Error(errCount)
		c.Response(500, "", nil)
	}
	Data.Total = cnt

	// 获取记录
	_, Err = qs.Limit(ArticleListSearchS.Limit).Offset(offset).All(&Data.List)
	if Err != nil {
		logs.Error(Err)
		c.Response(500, "", nil)
	}
	logs.Info(Data)
	c.Response(200, "", Data)
}

//发布新文章
func (c *ArticleController) Release() {
	var (
		ArticleModel = new(cms.ArticleModel)
	)
	_ = c.GetRequestJson(&ArticleModel, true)
	logs.Info(ArticleModel)

	//验证表单
	valid := validation.Validation{}
	valid.Required(ArticleModel.Title, "title")
	valid.Required(ArticleModel.CategoryId, "CategoryId")
	valid.Required(ArticleModel.Content, "Content")

	if valid.HasErrors() {
		// 如果有错误信息，证明验证没通过
		// 打印错误信息
		for _, err := range valid.Errors {
			errMsg := err.Key + err.Message
			logs.Info(errMsg)
			c.Response(304, errMsg, nil)
			break
		}
	}

	//确认栏目ID是否存在
	cateOrm := orm.NewOrm()
	cate := cms.CategoryModel{Id:ArticleModel.CategoryId}
	cateReadErr := cateOrm.Read(&cate)

	if cateReadErr == orm.ErrNoRows {
		logs.Info("查询不到")
		c.Response(601, "", nil)
	} else if cateReadErr != nil {
		logs.Info(cateReadErr.Error())
		c.Response(500, cateReadErr.Error(), nil)
	}

	//初始化一些数据
	ArticleModel.Status = 1                           //文章状态
	ArticleModel.CreateTime = time.Now().Unix()       //发布时间
	ArticleModel.UpdateTime = ArticleModel.CreateTime //更新时间
	ArticleModel.PostHits = 0                         //查看数
	ArticleModel.PostLike = 0                         //点赞数
	ArticleModel.CommentCount = 0                     //评论数

	//开启评论
	if ArticleModel.CommentStatus != 0 {
		ArticleModel.CommentStatus = 1
	} else {
		ArticleModel.CommentStatus = 0
	}

	//写入数据
	o := orm.NewOrm()
	id, err := o.Insert(ArticleModel)
	if err != nil || id <= 0 {
		logs.Error(err)
		c.Response(500, "", nil)
	}
	// 更新文章与栏目关联表
	Article := cms.ArticleModel{Id: id}
	_ = o.Read(&Article) // 查询最新的数据
	upErr := cms.UpdateCategoryArticles(Article, cate)
	if upErr != nil {
		logs.Error(upErr.Error())
		c.Response(500, "", nil)
	}
	c.Response(200, "", nil)
}

// 修改文章
func (c *ArticleController) Modify() {
	id, getErr := c.GetInt64("id")
	if getErr != nil {
		logs.Error(getErr.Error())
		c.Response(500, getErr.Error(), nil)
	}
	var (
		ArticleForm = new(cms.ArticleModel)
	)
	_ = c.GetRequestJson(&ArticleForm, true)
	logs.Info(ArticleForm)

	//验证表单
	if id == 0 {
		c.Response(304, "ID missing", nil)
	}
	valid := validation.Validation{}
	valid.Required(ArticleForm.Title, "title")
	valid.Required(ArticleForm.CategoryId, "CategoryId")
	valid.Required(ArticleForm.Content, "Content")

	if valid.HasErrors() {
		// 如果有错误信息，证明验证没通过
		// 打印错误信息
		for _, err := range valid.Errors {
			errMsg := err.Key + err.Message
			logs.Info(errMsg)
			c.Response(304, errMsg, nil)
			break
		}
	}

	o := orm.NewOrm()
	// 查找文章
	Article := cms.ArticleModel{Id: id}
	err := o.Read(&Article)

	if err == orm.ErrNoRows {
		logs.Error("查询不到")
		c.Response(602, "", nil)
	} else if err == orm.ErrMissPK {
		logs.Error("找不到主键")
		c.Response(401, "", nil)
	}

	//确认栏目ID是否存在
	cateOrm := orm.NewOrm()
	cate := cms.CategoryModel{Id:ArticleForm.CategoryId}
	cateReadErr := cateOrm.Read(&cate)

	if cateReadErr == orm.ErrNoRows {
		logs.Info("查询不到")
		c.Response(601, "", nil)
	} else if cateReadErr != nil {
		logs.Info(cateReadErr.Error())
		c.Response(500, cateReadErr.Error(), nil)
	}

	//初始化一些数据,保持这些数据不被人为修改
	ArticleForm.UpdateTime = time.Now().Unix()      //更新时间
	ArticleForm.PostHits = Article.PostHits         //查看数
	ArticleForm.PostLike = Article.PostLike         //点赞数
	ArticleForm.CommentCount = Article.CommentCount //评论数
	ArticleForm.Id = id

	//保存数据
	UpdateNum, UpdateErr := o.Update(ArticleForm)
	if UpdateErr != nil {
		logs.Error(err)
		c.Response(500, "", nil)
	}
	if UpdateNum > 0 {
		// 更新文章与栏目关联表
		Article = cms.ArticleModel{Id: id}
		_ = o.Read(&Article) // 查询最新的数据
		upErr := cms.UpdateCategoryArticles(Article, cate)
		if upErr != nil {
			logs.Error(upErr.Error())
			c.Response(500, "", nil)
		}
	} else {
		logs.Error("没有文章被更新")
		c.Response(603, "", nil)
	}
	c.Response(200, "", nil)
}
