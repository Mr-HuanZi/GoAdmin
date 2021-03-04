package cms

import (
	"github.com/beego/beego/v2/client/orm"
	"github.com/beego/beego/v2/core/logs"
	"go-admin/controllers/admin"
	"go-admin/models/cms"
	"go-admin/utils"
	"strconv"
	"time"
)

// 文章栏目
type CategoryController struct {
	admin.BaseController
}

type CategoryListResult struct {
	Total int64
	List  []*cms.CategoryModel
}

// 栏目列表
func (c *CategoryController) List() {
	var (
		CategoryListSearchS = &struct {
			Name  string
			Limit int `valid:"Range(0, 1000)"` //分页每页显示的条数
			Page  int `valid:"Min(1)"`         //当前页码
		}{}
		Category = new(cms.CategoryModel)
		Err      error
		Data     = new(CategoryListResult)
		offset   int
	)
	_ = c.GetRequestJson(&CategoryListSearchS, false)

	//获取每页记录条数, 页码, 计算页码偏移量
	CategoryListSearchS.Limit, CategoryListSearchS.Page, offset = c.Paginate(CategoryListSearchS.Page, CategoryListSearchS.Limit)

	o := orm.NewOrm()
	qs := o.QueryTable(Category)

	if CategoryListSearchS.Name != "" {
		qs = qs.Filter("name__contains", CategoryListSearchS.Name)
	}

	// 获取总条目
	cnt, errCount := qs.Count()
	if errCount != nil {
		logs.Error(errCount)
		c.Response(500, "", nil)
	}
	Data.Total = cnt

	// 获取记录
	_, Err = qs.Limit(CategoryListSearchS.Limit).Offset(offset).All(&Data.List)
	if Err != nil {
		logs.Error(Err)
		c.Response(500, "", nil)
	}
	logs.Info(Data)
	c.Response(200, "", Data)
}

// 新增栏目
func (c *CategoryController) Add() {
	var (
		CategoryForm = new(cms.CategoryModel)
		validateMsg  string
		validateRes  bool
	)
	_ = c.GetRequestJson(&CategoryForm, true)

	/* 表单字段验证 Start */
	validateRes, validateMsg = utils.FormValidation(CategoryForm)
	if !validateRes {
		c.Response(304, validateMsg, nil)
	}
	/* 表单字段验证 End */

	//初始化一些数据
	if CategoryForm.Status == 0 {
		CategoryForm.Status = 1
	}
	CategoryForm.CreateTime = time.Now().Unix()

	o := orm.NewOrm()
	//查找相同的别名
	cnt, errCount := o.QueryTable(new(cms.CategoryModel)).Filter("alias", CategoryForm.Alias).Count()
	if errCount != nil {
		logs.Error(errCount)
		c.Response(500, "", nil)
	}
	if cnt > 0 {
		//有重复
		c.Response(600, "", nil)
	}
	c.Response(200, "", nil)
	return
	_, errInsert := o.Insert(CategoryForm)
	if errInsert != nil {
		logs.Error(errInsert)
		c.Response(500, "", nil)
	}
	c.Response(200, "", nil)
}

// 修改栏目
func (c *CategoryController) Modify() {
	id, getErr := c.GetInt("id")
	if getErr != nil {
		c.Response(500, getErr.Error(), nil)
	}
	var (
		CategoryForm = new(cms.CategoryModel)
		validateMsg  string
		validateRes  bool
	)
	_ = c.GetRequestJson(&CategoryForm, true)

	/* 表单字段验证 Start */
	if id == 0 {
		c.Response(303, "", nil)
	}
	validateRes, validateMsg = utils.FormValidation(CategoryForm)
	if !validateRes {
		c.Response(304, validateMsg, nil)
	}
	/* 表单字段验证 End */
	CategoryForm.Id = id
	num, err := cms.UpdateCategory(CategoryForm)

	if err != nil {
		c.Response(500, "", nil)
	} else {
		c.Response(200, "", num)
	}
}

// 栏目删除
func (c *CategoryController) Delete() {
	id, getErr := c.GetInt("id")
	if getErr != nil {
		logs.Error(getErr.Error())
		c.Response(500, getErr.Error(), nil)
	}

	o := orm.NewOrm()
	if num, err := o.Delete(&cms.CategoryModel{Id: id}); err == nil {
		if num > 0 {
			c.Response(200, "", nil)
		} else {
			c.Response(405, "", nil)
		}
	} else {
		c.Response(500, err.Error(), nil)
	}
}

// 获取单个栏目信息
func (c *CategoryController) GetCategory() {
	strv := c.Ctx.Input.Param(":id")
	if len(strv) <= 0 {
		c.Response(500, "", nil)
	}
	id, getErr := strconv.Atoi(strv)
	if getErr != nil {
		logs.Error(getErr.Error())
		c.Response(500, getErr.Error(), nil)
	}

	// 查询栏目
	o := orm.NewOrm()
	Category := cms.CategoryModel{Id: id}
	err := o.Read(&Category)
	if err == orm.ErrNoRows {
		logs.Error("查询不到")
		c.Response(602, "", nil)
	} else if err == orm.ErrMissPK {
		logs.Error("找不到主键")
		c.Response(401, "", nil)
	} else {
		c.Response(200, "", Category)
	}

}
