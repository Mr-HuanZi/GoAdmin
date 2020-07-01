package admin

import (
	"github.com/astaxie/beego/logs"
	"github.com/astaxie/beego/orm"
	"go-admin/models/admin"
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
