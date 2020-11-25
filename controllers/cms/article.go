package cms

import (
	"github.com/astaxie/beego/logs"
	"github.com/astaxie/beego/orm"
	"github.com/astaxie/beego/validation"
	"go-admin/controllers/admin"
	"go-admin/lib"
	"go-admin/models/cms"
	"html"
	"time"
)

type ArticleController struct {
	admin.BaseController
}

func (c *ArticleController) List() {
	var (
		ArticleListSearchS = &struct {
			Limit      int       `valid:"Range(0, 1000)"` //分页每页显示的条数
			Page       int       `valid:"Min(1)"`         //当前页码
			Title      string    //文章标题
			CreateTime [2]string //文章创建时间
			Status     uint8     //文章状态
			Category   int       //文章分类
		}{}
		Data = &struct {
			Total int64
			List  []*cms.ArticleModel
		}{}
		Err     error
		Article = new(cms.ArticleModel)
		offset  int
	)
	_ = c.GetRequestJson(&ArticleListSearchS, false)
	logs.Info(ArticleListSearchS)

	//获取每页记录条数, 页码, 计算页码偏移量
	ArticleListSearchS.Limit, ArticleListSearchS.Page, offset = c.Paginate(ArticleListSearchS.Page, ArticleListSearchS.Limit)

	o := orm.NewOrm()
	qs := o.QueryTable(Article)

	//状态搜索
	if ArticleListSearchS.Status != 0 {
		qs = qs.Filter("status", ArticleListSearchS.Status)
	} else {
		qs = qs.Filter("status__in", 1, 2, 3)
	}

	// 开始时间
	if ArticleListSearchS.CreateTime[0] != "" {
		stamp, _ := time.ParseInLocation("2006-01-02 15:04:05", ArticleListSearchS.CreateTime[0], time.Local)
		qs = qs.Filter("create_time__gte", stamp.UnixNano()/1e6)
	}

	//结束时间
	if ArticleListSearchS.CreateTime[1] != "" {
		stamp, _ := time.ParseInLocation("2006-01-02 15:04:05", ArticleListSearchS.CreateTime[1], time.Local)
		qs = qs.Filter("create_time__lte", stamp.UnixNano()/1e6)
	}

	// 标题搜索
	if ArticleListSearchS.Title != "" {
		qs = qs.Filter("title__contains", ArticleListSearchS.Title)
	}

	// 获取总条目
	cnt, errCount := qs.Count()
	if errCount != nil {
		logs.Error(errCount)
		c.Response(500, "", nil)
	}
	Data.Total = cnt

	// 获取记录
	_, Err = qs.Limit(ArticleListSearchS.Limit).Offset(offset).OrderBy("-id").All(&Data.List)
	if Err != nil {
		logs.Error(Err)
		c.Response(500, "", nil)
	}

	// 获取栏目名称
	CategoryList := make(map[int]*cms.CategoryModel) // 缓存，避免单次的重复查询
	for i, item := range Data.List {
		if item.CategoryId != 0 {
			// 先从内存中查找数据
			if cate, ok := CategoryList[item.CategoryId]; ok {
				Data.List[i].CategoryName = cate.Name
			} else {
				cate := cms.CategoryModel{Id: item.CategoryId}
				err := o.Read(&cate)
				if err == orm.ErrNoRows {
					// 没有相应记录
					c.Response(404, "", nil)
					return
				} else if err == orm.ErrMissPK {
					// 主键丢失
					logs.Error(orm.ErrMissPK)
					c.Response(401, "", nil)
					return
				} else {
					CategoryList[item.CategoryId] = &cate //赋值到map，避免重复查询
					Data.List[i].CategoryName = cate.Name
				}
			}
		}
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

	// 对内容进行转义
	ArticleModel.Content = html.EscapeString(ArticleModel.Content)

	//验证表单
	valid := validation.Validation{}
	valid.Required(ArticleModel.Title, "title")
	valid.Required(ArticleModel.CategoryId, "category_id")
	valid.Required(ArticleModel.Content, "content")

	if valid.HasErrors() {
		// 如果有错误信息，证明验证没通过
		// 打印错误信息
		for _, err := range valid.Errors {
			errMsg := err.Key + " " + err.Message
			logs.Info(errMsg)
			c.Response(304, errMsg, nil)
			break
		}
	}

	//确认栏目ID是否存在
	cateOrm := orm.NewOrm()
	cate := cms.CategoryModel{Id: ArticleModel.CategoryId}
	cateReadErr := cateOrm.Read(&cate)

	if cateReadErr == orm.ErrNoRows {
		logs.Info("查询不到")
		c.Response(601, "", nil)
	} else if cateReadErr != nil {
		logs.Info(cateReadErr.Error())
		c.Response(500, cateReadErr.Error(), nil)
	}

	//初始化一些数据
	ArticleModel.Status = 1                               //文章状态
	ArticleModel.CreateTime = time.Now().UnixNano() / 1e6 //创建时间
	ArticleModel.UpdateTime = ArticleModel.CreateTime     //更新时间
	ArticleModel.PostHits = 0                             //查看数
	ArticleModel.PostLike = 0                             //点赞数
	ArticleModel.CommentCount = 0                         //评论数
	ArticleModel.Author = lib.CurrentUser.UserNickname    //作者
	ArticleModel.StaffId = lib.CurrentUser.Id             //作者ID

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

	// 对内容进行转义
	ArticleForm.Content = html.EscapeString(ArticleForm.Content)

	//验证表单
	if id == 0 {
		c.Response(303, "", nil)
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
	cate := cms.CategoryModel{Id: ArticleForm.CategoryId}
	cateReadErr := cateOrm.Read(&cate)

	if cateReadErr == orm.ErrNoRows {
		logs.Info("查询不到")
		c.Response(601, "", nil)
	} else if cateReadErr != nil {
		logs.Info(cateReadErr.Error())
		c.Response(500, cateReadErr.Error(), nil)
	}

	//初始化一些数据,保持这些数据不被人为修改
	ArticleForm.UpdateTime = time.Now().UnixNano() / 1e6 //更新时间
	ArticleForm.PostHits = Article.PostHits              //查看数
	ArticleForm.PostLike = Article.PostLike              //点赞数
	ArticleForm.CommentCount = Article.CommentCount      //评论数
	ArticleForm.Author = lib.CurrentUser.UserNickname    //作者
	ArticleForm.StaffId = lib.CurrentUser.Id             //作者ID
	ArticleForm.Status = Article.Status                  //作者ID
	ArticleForm.Id = id

	//保存数据
	UpdateNum, UpdateErr := o.Update(ArticleForm)
	if UpdateErr != nil {
		logs.Error(UpdateErr)
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
		c.Response(403, "", nil)
	}
	c.Response(200, "", nil)
}

// 文章删除
func (c *ArticleController) Delete() {
	id, getErr := c.GetInt64("id")
	if getErr != nil {
		logs.Error(getErr.Error())
		c.Response(500, getErr.Error(), nil)
	}

	c.Response(200, "", nil)
	return

	o := orm.NewOrm()
	if num, err := o.Delete(&cms.ArticleModel{Id: id}); err == nil {
		if num > 0 {
			// 删除关联表
			_, delErr := o.Delete(&cms.CategoryArticlesModel{ArticlesId: id})
			if delErr != nil {
				// 如果出错了也不中断
				logs.Error(delErr)
			}
			c.Response(200, "", nil)
		} else {
			c.Response(405, "", nil)
		}
	} else {
		c.Response(500, err.Error(), nil)
	}
}

// 获取一篇文章信息
func (c *ArticleController) GetArticle() {
	id, getErr := c.GetInt64("id")
	if getErr != nil {
		logs.Error(getErr.Error())
		c.Response(500, getErr.Error(), nil)
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
	} else {
		Article.Content = html.UnescapeString(Article.Content)
		c.Response(200, "", Article)
	}
}
