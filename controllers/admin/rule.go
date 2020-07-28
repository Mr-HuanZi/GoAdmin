package admin

import (
	"fmt"
	"github.com/astaxie/beego/logs"
	"github.com/astaxie/beego/orm"
	"go-admin/lib"
	"go-admin/models/admin"
	"reflect"
	"strconv"
	"strings"
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

	o := orm.NewOrm()
	/* 表单字段验证 Start */
	validateRes, validateMsg = lib.FormValidation(RuleForm)
	if !validateRes {
		c.Response(304, validateMsg, nil)
	}
	// 检查重复项
	ruleCount, ruleCountErr := o.QueryTable(new(admin.RuleModel)).Filter("rule", RuleForm.Rule).Count()
	if ruleCountErr != nil {
		logs.Error(ruleCountErr)
		c.Response(500, ruleCountErr.Error(), nil)
	}
	if ruleCount >= 1 {
		c.Response(402, "", nil)
	}
	/* 表单字段验证 End */

	// 初始化一些字段
	RuleForm.CreateTime = time.Now().Unix()
	RuleForm.AddStaff = c.ThatUser.Id
	RuleForm.Id = 0

	// 写入数据
	_, err := o.Insert(&RuleForm)
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

	o := orm.NewOrm()
	/* 表单字段验证 Start */
	if id == 0 {
		c.Response(304, "ID missing", nil)
	}
	validateRes, validateMsg = lib.FormValidation(RuleForm)
	if !validateRes {
		c.Response(304, validateMsg, nil)
	}
	// 检查相同的记录
	ruleCount, ruleCountErr := o.QueryTable(new(admin.RuleModel)).Filter("rule", RuleForm.Rule).Exclude("id", id).Count()
	if ruleCountErr != nil {
		logs.Error(ruleCountErr)
		c.Response(500, ruleCountErr.Error(), nil)
	}
	if ruleCount >= 1 {
		c.Response(402, "", nil)
	}
	/* 表单字段验证 End */
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
	RuleForm.Id = Rule.Id // 保持ID不改变

	//保存数据
	UpdateNum, UpdateErr := o.Update(&RuleForm)
	if UpdateErr != nil {
		logs.Error(err)
		c.Response(500, "", nil)
	}
	fmt.Println(UpdateNum)
	c.Response(200, "", nil)
}

// 写入权限组数据
func (c *RuleController) WriteGroup() {
	var (
		AuthGroup admin.AuthGroupModel
		id        int
		AtoiErr   error
	)
	// 直接用GetInt的话，如果id不存在，会报错
	queryS := c.Ctx.Input.Query("id")
	if len(queryS) > 0 {
		id, AtoiErr = strconv.Atoi(queryS)
		if AtoiErr != nil {
			logs.Error(AtoiErr.Error())
			c.Response(500, AtoiErr.Error(), nil)
		}
	}

	// 结构体
	groupForm := &struct {
		Title       string `valid:"Required"`
		Description string
	}{}

	_ = c.GetRequestJson(&groupForm, true)

	/* 表单字段验证 Start */
	validateRes, validateMsg := lib.FormValidation(groupForm)
	if !validateRes {
		c.Response(304, validateMsg, nil)
	}
	/* 表单字段验证 End */

	o := orm.NewOrm()
	if id > 0 {
		AuthGroup.Id = id
		readErr := o.Read(&AuthGroup)
		if readErr == orm.ErrNoRows {
			logs.Info("没有相关记录")
			c.Response(404, "", nil)
		} else if readErr == orm.ErrMissPK {
			logs.Info("找不到主键")
			c.Response(401, "", nil)
		}
	}
	// 写数据到结构体里
	AuthGroup.Title = groupForm.Title
	AuthGroup.Description = groupForm.Description
	if AuthGroup.Id > 0 {
		// 更新
		num, err := o.Update(&AuthGroup)
		if err != nil {
			logs.Error(err.Error())
			c.Response(500, "", nil)
		}
		if num <= 0 {
			logs.Info("更新记录为0")
			c.Response(403, "", nil)
		}
	} else {
		// 新增
		AuthGroup.Status = 0
		_, insertErr := o.Insert(&AuthGroup)
		if insertErr != nil {
			logs.Error(insertErr.Error())
			c.Response(500, "", nil)
		}
	}
	c.Response(200, "", nil)
}

// 权限组授权
func (c *RuleController) AccessAuth() {
	id, getErr := c.GetInt("id")
	if getErr != nil {
		logs.Error(getErr.Error())
		c.Response(500, "", nil)
	}

	// 查询组是否存在
	AuthGroup := admin.AuthGroupModel{Id: id}
	o := orm.NewOrm()
	readErr := o.Read(&AuthGroup)
	if readErr == orm.ErrNoRows {
		logs.Info("没有相关记录")
		c.Response(404, "", nil)
	} else if readErr == orm.ErrMissPK {
		logs.Info("找不到主键")
		c.Response(401, "", nil)
	}

	// 获取请求数据
	rulesJson := &struct {
		Rules string `valid:"Required"`
	}{}

	_ = c.GetRequestJson(&rulesJson, true)

	/* 表单字段验证 Start */
	validateRes, validateMsg := lib.FormValidation(rulesJson)
	if !validateRes {
		c.Response(304, validateMsg, nil)
	}
	/* 表单字段验证 End */

	// 尝试分割字符串
	rules := strings.Split(rulesJson.Rules, ",")
	if len(rules) <= 0 {
		logs.Info("字符串分割结果为空")
		c.Response(502, "", nil)
	}
	fmt.Println(rules)
	// 去除空元素
	var rulesResult []string
	for _, value := range rules {
		if value != "" {
			rule, err := strconv.Atoi(value)
			if err == nil {
				rulesResult = append(rulesResult, strconv.Itoa(rule)) //把变量转换回字符串
			}
		}
	}
	fmt.Println(rulesResult)

	// 合并字符串
	AuthGroup.Rules = strings.Join(rulesResult, ",")
	// 更新数据
	num, updateErr := o.Update(&AuthGroup, "rules")
	if updateErr != nil {
		logs.Error(updateErr.Error())
		c.Response(500, updateErr.Error(), nil)
	}
	if num <= 0 {
		logs.Info("没有记录被更新")
		c.Response(403, "", nil)
	}
	c.Response(200, "", nil)
}

// 人员授权才做
func (c *RuleController) MemberAuth() {
	id, getErr := c.GetInt("id")
	if getErr != nil {
		logs.Error(getErr.Error())
		c.Response(500, "", nil)
	}
	// 查询组是否存在
	AuthGroup := admin.AuthGroupModel{Id: id}
	o := orm.NewOrm()
	readErr := o.Read(&AuthGroup)
	if readErr == orm.ErrNoRows {
		logs.Info("没有相关记录")
		c.Response(404, "", nil)
	} else if readErr == orm.ErrMissPK {
		logs.Info("找不到主键")
		c.Response(401, "", nil)
	}

	postJson := &struct {
		Uid []interface{} `valid:"Required"`
	}{}

	_ = c.GetRequestJson(&postJson, true)

	/* 表单字段验证 Start */
	validateRes, validateMsg := lib.FormValidation(postJson)
	if !validateRes {
		c.Response(304, validateMsg, nil)
	}
	/* 表单字段验证 End */

	actType := c.GetString("type")
	if actType == "auth" {
		// 成员授权
		for _, v := range postJson.Uid {
			if reflect.TypeOf(v).String() != "float64" {
				// 如果类型不对，则终止
				c.Response(503, "", nil)
				break
			}
		}
	}
}
