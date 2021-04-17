package controllers

import (
	"github.com/beego/beego/v2/core/logs"
	"go-admin/models"
	"go-admin/utils"
	"time"
)

type LoginController struct {
	BaseController
}

//登录表单
type LoginFormS struct {
	Username string `valid:"Required;MinSize(4);MaxSize(18);AlphaDash"`
	Password string `valid:"Required;MinSize(6);MaxSize(18);AlphaDash"`
}

//注册表单
type RegisterFormS struct {
	Username   string `valid:"Required;MinSize(4);MaxSize(18);AlphaDash"`
	Password   string `valid:"Required;MinSize(6);MaxSize(18);AlphaDash"`
	RePassword string `valid:"Required;MinSize(6);MaxSize(18);AlphaDash"`
	Email      string `valid:"Required;Email"`
}

//登录
func (c *LoginController) Login() {
	var (
		loginForm LoginFormS
		code      int
		uid       int64
	)
	_ = c.GetRequestJson(&loginForm, true)
	/* 表单字段验证 Start */
	c.FormValidation(&loginForm)
	/* 表单字段验证 End */
	//对密码加密
	loginForm.Password = utils.Encryption(loginForm.Password)
	//验证用户登录
	code, uid = models.Login(loginForm.Username, loginForm.Password)
	if code == 100 {
		//登录成功
		token, tokenErr := utils.GenerateUserToken(uid) //获取登录令牌
		if tokenErr != nil {
			c.Response(101, "", nil) //令牌生成失败
		}
		// 获取用户信息
		_ = c.getLoginUser(uid)
		//记录用户登录信息
		models.UpdateUserLoginInfo(uid, c.Ctx.Input.IP())
		c.Ctx.SetCookie("Authorization", token, 7200, "/", "", false, true)
		c.Ctx.Output.Header("Authorization", token)
	}
	c.Response(code, "", nil)
}

func (c *LoginController) Test() {
	logs.Info("Test is run.")
}

//注册用户
func (c *LoginController) Register() {
	var (
		registerForm RegisterFormS
	)
	_ = c.GetRequestJson(&registerForm, true)
	//验证表单
	c.FormValidation(&registerForm)
	//确认密码是否相等
	if registerForm.Password != registerForm.RePassword {
		c.Response(105, "", nil) //设置的密码与确认密码不一致
	}
	//检查用户名和邮箱是否已被注册
	if models.CheckUserRepeat(registerForm.Username, registerForm.Email) {
		c.Response(107, "", nil) //用户名或邮箱已被注册
	}
	//创建用户数据
	var UserData models.UserModel
	var uid int64
	UserData.UserLogin = registerForm.Username
	UserData.UserPass = utils.Encryption(registerForm.Password) //加密密码
	UserData.UserNickname = registerForm.Username               //用户昵称默认是登录账号
	UserData.UserType = 1                                       //管理员类型
	UserData.CreateTime = time.Now().Unix()                     //管理员类型
	UserData.UpdateTime = UserData.CreateTime                   //管理员类型
	UserData.UserStatus = 1                                     //用户状态
	uid = models.CreateUser(&UserData)
	if uid != 0 {
		logs.Info(uid)
		c.Response(110, "", nil) //注册成功
	}
	c.Response(400, "", nil) //注册失败
}
