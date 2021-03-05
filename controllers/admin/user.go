package admin

import (
	"github.com/beego/beego/v2/client/orm"
	"github.com/beego/beego/v2/core/logs"
	"go-admin/models/admin"
	"go-admin/utils"
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

	//获取每页记录条数, 页码, 计算页码偏移量
	SearchFormInstance.Limit, SearchFormInstance.Page, offset = c.Paginate(SearchFormInstance.Page, SearchFormInstance.Limit)

	o := orm.NewOrm()
	qs := o.QueryTable(UserModel)

	// 用户名搜索
	if SearchFormInstance.Username != "" {
		qs = qs.Filter("user_login__contains", SearchFormInstance.Username)
	}

	// 用户昵称搜索
	if SearchFormInstance.Nickname != "" {
		qs = qs.Filter("user_nickname__contains", SearchFormInstance.Username)
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
		UserForm admin.UserModel
	)

	_ = c.GetRequestJson(&UserForm, true)

	/* 表单字段验证 Start */
	c.FormValidation(UserForm)
	/* 表单字段验证 End */

	// 初始化一些字段
	UserForm.CreateTime = time.Now().Unix()
	UserForm.UpdateTime = UserForm.CreateTime
	UserForm.First = 1
	UserForm.UserStatus = 1
	UserForm.UserType = 1 //管理员类型
	// 密码加密
	UserForm.UserPass = utils.Encryption(UserForm.UserPass)

	// 检查用户名是否存在
	duplication := admin.CheckUserDuplication(UserForm.UserLogin)
	if duplication {
		c.Response(107, "", nil)
	}
	// 写入数据
	o := orm.NewOrm()
	_, err := o.Insert(&UserForm)
	if err != nil {
		logs.Error(err)
		c.Response(500, "", nil)
	}
	c.Response(200, "", nil)

}

// 修改用户信息
func (c *UserController) Modify() {
	id, getErr := c.GetInt64("id")
	if getErr != nil {
		logs.Error(getErr.Error())
		c.Response(500, getErr.Error(), nil)
	}

	var (
		UserForm admin.UserModel
	)

	_ = c.GetRequestJson(&UserForm, true)

	/* 表单字段验证 Start */
	if id == 0 {
		c.Response(304, "ID missing", nil)
	}
	c.FormValidation(UserForm)
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

	// 密码加密
	UserForm.UserPass = utils.Encryption(UserForm.UserPass)

	// 不能被修改的数据
	UserForm.CreateTime = User.CreateTime
	UserForm.UpdateTime = User.UpdateTime
	UserForm.First = User.First
	UserForm.UserStatus = User.UserStatus
	UserForm.UserType = 1 //管理员类型
	UserForm.Id = User.Id

	// 如果修改了用户名，检查用户名是否存在
	if UserForm.UserLogin != User.UserLogin {
		duplication := admin.CheckUserDuplication(UserForm.UserLogin)
		if duplication {
			c.Response(107, "", nil)
		}
	}

	//保存数据
	UpdateNum, UpdateErr := o.Update(&UserForm)
	if UpdateErr != nil {
		logs.Error(err)
		c.Response(500, "", nil)
	}
	if UpdateNum <= 0 {
		logs.Notice("更新用户信息的记录为[", UpdateNum, "]")
		c.Response(403, "", nil)
	}
	c.Response(200, "", nil)
}

// 禁用用户
func (c *UserController) ForbidUser() {
	c.changeUserStatus(0)
}

// 启用用户
func (c *UserController) ResumeUser() {
	c.changeUserStatus(1)
}

// 更改用户状态
func (c *UserController) changeUserStatus(status int8) {
	var UidMap map[string][]int64
	_ = c.GetRequestJson(&UidMap, true)

	/* 表单字段验证 Start */
	if _, ok := UidMap["Uid"]; !ok {
		c.Response(304, "找不到Uid字段", nil)
	}
	if len(UidMap["Uid"]) <= 0 {
		c.Response(304, "Uid不能为空", nil)
	}
	/* 表单字段验证 End */

	updateNum, updateErr := admin.ChangeUserStatus(UidMap["Uid"], status)
	if updateErr != nil {
		c.Response(500, "", nil)
	}
	if updateNum <= 0 {
		c.Response(405, "", nil)
	}
	c.Response(200, "", nil)
}

// 获取当前用户信息
func (c *UserController) FetchCurrentUser() {
	if utils.CurrentUser.Id != 0 {
		user := make(map[string]interface{})
		user["userId"] = utils.CurrentUser.Id
		user["nickname"] = utils.CurrentUser.UserNickname
		user["username"] = utils.CurrentUser.UserLogin
		user["email"] = utils.CurrentUser.UserEmail
		user["birthday"] = utils.CurrentUser.Birthday
		user["mobile"] = utils.CurrentUser.Mobile
		if utils.CurrentUser.Sex == 1 {
			user["sex"] = "男"
		} else if utils.CurrentUser.Sex == 2 {
			user["sex"] = "女"
		} else {
			user["sex"] = "保密"
		}
		user["signature"] = utils.CurrentUser.Signature
		user["user_url"] = utils.CurrentUser.UserUrl
		user["avatar"] = utils.CurrentUser.Avatar
		c.Response(200, "", user)
	} else {
		c.Response(200, "", nil)
	}
}
