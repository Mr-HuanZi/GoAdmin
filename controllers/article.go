package controllers

import (
	"github.com/beego/beego/v2/client/orm"
	"github.com/beego/beego/v2/core/logs"
	"github.com/beego/beego/v2/core/validation"
	"go-admin/models"
	"go-admin/utils"
	"html"
	"time"
)

const (
	TOP       = 0b001
	RECOMMEND = 0b010
)

type ArticleController struct {
	BaseController
}

type ArticleFormS struct {
	Title         string          `json:"title"`
	CategoryId    int             `json:"category_id,string"`
	Describe      string          `json:"describe"`
	Content       string          `json:"content"`
	Status        int8            `json:"status"`
	CreateTime    int64           `json:"create_time"`
	UpdateTime    int64           `json:"update_time"`
	Tag           string          `json:"tag"`
	PostHits      int64           `json:"post_hits,string"`
	PostLike      int64           `json:"post_like,string"`
	CommentCount  int64           `json:"comment_count,string"`
	CommentStatus int             `json:"comment_status"`
	More          string          `json:"more"`
	Source        string          `json:"source"`
	Recommend     map[string]bool `json:"recommend"`
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
			List  []*models.ArticleModel
		}{}
		Err     error
		Article = new(models.ArticleModel)
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
	}

	// 开始时间
	if ArticleListSearchS.CreateTime[0] != "" {
		stamp, _ := time.ParseInLocation(utils.DatetimeFormat, ArticleListSearchS.CreateTime[0], time.Local)
		qs = qs.Filter("create_time__gte", stamp.UnixNano()/1e6)
	}

	//结束时间
	if ArticleListSearchS.CreateTime[1] != "" {
		stamp, _ := time.ParseInLocation(utils.DatetimeFormat, ArticleListSearchS.CreateTime[1], time.Local)
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
	CategoryList := make(map[int]*models.CategoryModel) // 缓存，避免单次的重复查询
	for i, item := range Data.List {
		if item.CategoryId != 0 {
			// 先从内存中查找数据
			if cate, ok := CategoryList[item.CategoryId]; ok {
				Data.List[i].CategoryName = cate.Name
			} else {
				cate := models.CategoryModel{Id: item.CategoryId}
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
		ArticleModel = new(models.ArticleModel)
		ArticleForm  = new(ArticleFormS)
	)
	_ = c.GetRequestJson(&ArticleForm, true)

	// 对内容进行转义
	ArticleForm.Content = html.EscapeString(ArticleForm.Content)

	//验证表单
	valid := validation.Validation{}
	valid.Required(ArticleForm.Title, "title")
	valid.Required(ArticleForm.CategoryId, "category_id")
	valid.Required(ArticleForm.Content, "content")
	valid.Range(ArticleForm.Status, 0, 2, "status")

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
	cate := models.CategoryModel{Id: ArticleForm.CategoryId}
	cateReadErr := cateOrm.Read(&cate)

	if cateReadErr == orm.ErrNoRows {
		logs.Info("查询不到")
		c.Response(601, "", nil)
	} else if cateReadErr != nil {
		logs.Info(cateReadErr.Error())
		c.Response(500, cateReadErr.Error(), nil)
	}

	// 把ArticleForm的值赋值给ArticleModel
	StructCopyErr := utils.StructCopy(ArticleModel, ArticleForm)
	if StructCopyErr != nil {
		logs.Error(StructCopyErr)
		c.Response(500, "", nil)
		return
	}

	//初始化一些数据
	ArticleModel.CreateTime = utils.UnixMilli()          //创建时间
	ArticleModel.UpdateTime = ArticleModel.CreateTime    //更新时间
	ArticleModel.PostHits = 0                            //查看数
	ArticleModel.PostLike = 0                            //点赞数
	ArticleModel.CommentCount = 0                        //评论数
	ArticleModel.Author = utils.CurrentUser.UserNickname //作者
	ArticleModel.StaffId = utils.CurrentUser.Id          //作者ID

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
	Article := models.ArticleModel{Id: id}
	_ = o.Read(&Article) // 查询最新的数据
	upErr := models.UpdateCategoryArticles(Article, cate)
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
		ArticleModel = new(models.ArticleModel)
		ArticleForm  = new(ArticleFormS)
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
	Article := models.ArticleModel{Id: id}
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
	cate := models.CategoryModel{Id: ArticleForm.CategoryId}
	cateReadErr := cateOrm.Read(&cate)

	if cateReadErr == orm.ErrNoRows {
		logs.Info("查询不到")
		c.Response(601, "", nil)
	} else if cateReadErr != nil {
		logs.Info(cateReadErr.Error())
		c.Response(500, "", nil)
	}

	// 计算状态
	flag := int(Article.Recommend)
	if _, ok := ArticleForm.Recommend["top"]; ok {
		if ArticleForm.Recommend["top"] {
			flag = utils.ShiftFlag(true, TOP, flag)
		} else {
			flag = utils.ShiftFlag(false, TOP, flag)
		}
	}
	if _, ok := ArticleForm.Recommend["recommend"]; ok {
		if ArticleForm.Recommend["recommend"] {
			flag = utils.ShiftFlag(true, RECOMMEND, flag)
		} else {
			flag = utils.ShiftFlag(false, RECOMMEND, flag)
		}
	}

	// 把ArticleForm的值赋值给ArticleModel
	StructCopyErr := utils.StructCopy(ArticleModel, ArticleForm)
	if StructCopyErr != nil {
		logs.Error(StructCopyErr)
		c.Response(500, "", nil)
		return
	}

	//初始化一些数据,保持这些数据不被人为修改
	ArticleModel.CreateTime = Article.CreateTime         //新增时间
	ArticleModel.UpdateTime = utils.UnixMilli()          //更新时间
	ArticleModel.Author = utils.CurrentUser.UserNickname //作者
	ArticleModel.StaffId = utils.CurrentUser.Id          //作者ID
	ArticleModel.Recommend = int8(flag)                  //推荐位
	ArticleModel.Id = id

	//保存数据
	UpdateNum, UpdateErr := o.Update(ArticleModel)
	if UpdateErr != nil {
		logs.Error(UpdateErr)
		c.Response(500, "", nil)
	}
	if UpdateNum > 0 {
		// 更新文章与栏目关联表
		Article = models.ArticleModel{Id: id}
		_ = o.Read(&Article) // 查询最新的数据
		upErr := models.UpdateCategoryArticles(Article, cate)
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

	o := orm.NewOrm()
	if num, err := o.Delete(&models.ArticleModel{Id: id}); err == nil {
		if num > 0 {
			// 删除关联表
			_, delErr := o.Delete(&models.CategoryArticlesModel{ArticlesId: id})
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
	Article := models.ArticleModel{Id: id}
	err := o.Read(&Article)

	if err == orm.ErrNoRows {
		logs.Error("查询不到")
		c.Response(602, "", nil)
	} else if err == orm.ErrMissPK {
		logs.Error("找不到主键")
		c.Response(401, "", nil)
	} else {
		Article.Content = html.UnescapeString(Article.Content)
		ArticleData := new(ArticleFormS)
		// 把Article的值赋值给ArticleData
		StructCopyErr := utils.StructCopy(ArticleData, &Article)
		if StructCopyErr != nil {
			logs.Error(StructCopyErr)
			c.Response(500, "", nil)
			return
		}

		// 推荐位
		ArticleData.Recommend = make(map[string]bool)
		ArticleData.Recommend["top"] = int(Article.Recommend)&TOP == TOP
		ArticleData.Recommend["recommend"] = int(Article.Recommend)&RECOMMEND == RECOMMEND
		c.Response(200, "", ArticleData)
	}
}
