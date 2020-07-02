package admin

import (
	"fmt"
	"github.com/astaxie/beego/logs"
	"github.com/astaxie/beego/orm"
	"go-admin/lib"
	"go-admin/models/admin"
	"time"
)

// 权限规则管理
type RuleController struct {
	BaseController
}

func (c *RuleController) List() {
	search := &struct {
		Limit  int `form:"limit"`
		Page   int `form:"page"`
		offset int
	}{}
	listResult := &struct {
		Total int64
		List  []*admin.RuleModel
	}{}
	if parseFormErr := c.ParseForm(search); parseFormErr != nil {
		logs.Error(parseFormErr.Error())
		c.Response(500, "", nil)
	}

	//获取每页记录条数, 页码, 计算页码偏移量
	search.Limit, search.Page, search.offset = c.Paginate(search.Page, search.Limit)

	o := orm.NewOrm()
	qs := o.QueryTable(new(admin.RuleModel))

	// 获取总条目
	cnt, errCount := qs.Count()
	if errCount != nil {
		logs.Error(errCount)
		c.Response(500, "", nil)
	}
	listResult.Total = cnt

	// 获取记录
	_, Err := qs.Limit(search.Limit).Offset(search.offset).All(&listResult.List)
	if Err != nil {
		logs.Error(Err)
		c.Response(500, "", nil)
	}
	logs.Info(listResult)
	c.Response(200, "", listResult)
}

func (c *RuleController) Add() {
	var (
		RuleForm    admin.RuleModel
		validateMsg string
		validateRes bool
	)

	_ = c.GetRequestJson(&RuleForm, true)
	logs.Info(RuleForm)

	/* 表单字段验证 Start */
	validateRes, validateMsg = lib.FormValidation(RuleForm)
	if !validateRes {
		c.Response(304, validateMsg, nil)
	}
	/* 表单字段验证 End */

	// 初始化一些字段
	RuleForm.CreateTime = time.Now().Unix()
	RuleForm.AddStaff = c.ThatUser.Id

	// 写入数据
	o := orm.NewOrm()
	_, err := o.Insert(RuleForm)
	if err != nil {
		logs.Error(err)
		c.Response(500, "", nil)
	}
	c.Response(200, "", nil)
}

func (c *RuleController) Modify() {
	id, getErr := c.GetInt("id")
	if getErr != nil {
		logs.Error(getErr.Error())
		c.Response(500, getErr.Error(), nil)
	}

	var (
		RuleForm    admin.RuleModel
		validateMsg string
		validateRes bool
	)

	_ = c.GetRequestJson(&RuleForm, true)
	logs.Info(RuleForm)

	/* 表单字段验证 Start */
	if id == 0 {
		c.Response(304, "ID missing", nil)
	}
	validateRes, validateMsg = lib.FormValidation(RuleForm)
	if !validateRes {
		c.Response(304, validateMsg, nil)
	}
	/* 表单字段验证 End */
	o := orm.NewOrm()
	// 查找文章
	Rule := admin.RuleModel{Id: id}
	err := o.Read(&Rule)

	if err == orm.ErrNoRows {
		logs.Error("查询不到")
		c.Response(602, "", nil)
	} else if err == orm.ErrMissPK {
		logs.Error("找不到主键")
		c.Response(401, "", nil)
	}

	// 不能被修改的数据
	RuleForm.CreateTime = Rule.CreateTime

	//保存数据
	UpdateNum, UpdateErr := o.Update(RuleForm)
	if UpdateErr != nil {
		logs.Error(err)
		c.Response(500, "", nil)
	}
	fmt.Println(UpdateNum)
	c.Response(200, "", nil)
}
