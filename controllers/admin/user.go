package admin

import (
	"fmt"
	"github.com/astaxie/beego/logs"
	"github.com/astaxie/beego/orm"
	"go-admin/lib"
	"go-admin/models/admin"
	"time"
)

// 后台管理员用户控制器
type UserController struct {
	BaseController
}

type SearchForm struct {
	Limit    int    `form:"limit"`
	Page     int    `form:"page"`
	Username string `form:"-"`
	Nickname string
}

type ListResult struct {
	Total int64
	List  []*admin.UserModel
}

// 用户列表
func (c *UserController) List() {
	var (
		SearchFormInstance SearchForm
		Data               = new(ListResult)
		UserModel          = new(admin.UserModel)
		Err                error
		offset             int
	)

	_ = c.GetRequestJson(&SearchFormInstance, false)
	logs.Info(SearchFormInstance)

	//获取每页记录条数, 页码, 计算页码偏移量
	SearchFormInstance.Limit, SearchFormInstance.Page, offset = c.Paginate(SearchFormInstance.Page, SearchFormInstance.Limit)

	o := orm.NewOrm()
	qs := o.QueryTable(UserModel)

	// 用户名搜索
	if SearchFormInstance.Username != "" {
		qs.Filter("user_login__contains", SearchFormInstance.Username)
	}

	// 用户昵称搜索
	if SearchFormInstance.Nickname != "" {
		qs.Filter("user_nickname__contains", SearchFormInstance.Username)
	}

	// 获取总条目
	cnt, errCount := qs.Count()
	if errCount != nil {
		logs.Error(errCount)
		c.Response(500, "", nil)
	}
	Data.Total = cnt

	// 获取记录
	_, Err = qs.Limit(SearchFormInstance.Limit).Offset(offset).All(&Data.List)
	if Err != nil {
		logs.Error(Err)
		c.Response(500, "", nil)
	}
	logs.Info(Data)
	c.Response(200, "", Data)
}

// 创建用户(管理员)
func (c *UserController) CreateUser() {
	var (
		UserForm    admin.UserModel
		validateMsg string
		validateRes bool
	)

	_ = c.GetRequestJson(&UserForm, true)
	logs.Info(UserForm)

	/* 表单字段验证 Start */
	validateRes, validateMsg = lib.FormValidation(UserForm)
	if !validateRes {
		c.Response(304, validateMsg, nil)
	}
	/* 表单字段验证 End */

	// 初始化一些字段
	UserForm.CreateTime = time.Now().Unix()
	UserForm.UpdateTime = UserForm.CreateTime
	UserForm.First = 1
	UserForm.UserStatus = 1
	UserForm.UserType = 1 //管理员类型

	// 写入数据
	o := orm.NewOrm()
	_, err := o.Insert(UserForm)
	if err != nil {
		logs.Error(err)
		c.Response(500, "", nil)
	}
	c.Response(200, "", nil)

}

func (c *UserController) Modify() {
	id, getErr := c.GetInt64("id")
	if getErr != nil {
		logs.Error(getErr.Error())
		c.Response(500, getErr.Error(), nil)
	}

	var (
		UserForm    admin.UserModel
		validateMsg string
		validateRes bool
	)

	_ = c.GetRequestJson(&UserForm, true)
	logs.Info(UserForm)

	/* 表单字段验证 Start */
	if id == 0 {
		c.Response(304, "ID missing", nil)
	}

	validateRes, validateMsg = lib.FormValidation(UserForm)
	if !validateRes {
		c.Response(304, validateMsg, nil)
	}
	/* 表单字段验证 End */
	o := orm.NewOrm()
	// 查找文章
	User := admin.UserModel{Id: id}
	err := o.Read(&User)

	if err == orm.ErrNoRows {
		logs.Error("查询不到")
		c.Response(602, "", nil)
	} else if err == orm.ErrMissPK {
		logs.Error("找不到主键")
		c.Response(401, "", nil)
	}

	// 不能被修改的数据
	UserForm.CreateTime = User.CreateTime
	UserForm.UpdateTime = User.UpdateTime
	UserForm.First = User.First
	UserForm.UserStatus = User.UserStatus
	UserForm.UserType = 1 //管理员类型

	//保存数据
	UpdateNum, UpdateErr := o.Update(UserForm)
	if UpdateErr != nil {
		logs.Error(err)
		c.Response(500, "", nil)
	}
	fmt.Println(UpdateNum)
	c.Response(200, "", nil)
}
