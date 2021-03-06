package controllers

import (
	"encoding/json"
	"errors"
	"github.com/beego/beego/v2/core/logs"
	"github.com/beego/beego/v2/core/validation"
	beego "github.com/beego/beego/v2/server/web"
	"go-admin/models"
	"go-admin/utils"
	"strconv"
)

type BaseController struct {
	beego.Controller
}

//实例化日志模块
var StatusCodeInstance utils.StatusCode

// 基类初始化方法，每次执行都会被调用
func (base *BaseController) Prepare() {
	//验证登录
	if !utils.InRuleExclude(base.Ctx.Request.RequestURI) {
		// 先从HTTP头中获取Authorization
		Authorization := base.Ctx.Input.Header("Authorization")
		if Authorization == "" {
			// 如果不存在再从cookie中获取
			Authorization = base.Ctx.GetCookie("Authorization")
		}
		// 第二次判断Authorization是否存在
		if Authorization == "" {
			// 直接返回登录失败
			base.Response(103, "", nil)
		}
		b, t := utils.ValidateUserToken(Authorization)
		if !b {
			//token验证失败
			base.Response(103, "", nil)
		} else {
			// 获取当前登录用户信息
			getUserErr := base.getLoginUser(t.User)
			if getUserErr != nil {
				base.Response(500, getUserErr.Error(), nil) //令牌生成失败
			}
			// 刷新令牌
			token, tokenErr := utils.RefreshUserToken(t)
			if tokenErr != nil {
				base.Response(101, "", nil) //令牌生成失败
			}
			if token != "" {
				base.Ctx.SetCookie("Authorization", token, 7200, "/", "", false, true)
				base.Ctx.Output.Header("Authorization", token)
			}
		}
		// 初始化权限规则
		base.initRule()
	}
}

// 执行完相应的 HTTP Method 方法之后执行
func (base *BaseController) Finish() {
}

func (base *BaseController) Response(code int, msg string, data interface{}) {
	statusCode := StatusCodeInstance.CreateData(code, msg, data)
	logs.Debug(statusCode)
	//if code == 103 {
	//	// 更改HTTP状态码
	//	base.Ctx.Output.Status = 401
	//}
	base.Data["json"] = &statusCode
	_ = base.ServeJSON()
	if code != 100 && code != 200 {
		base.StopRun()
	}
}

// 获取json请求数据
// @param s 接收数据的结构体
// @param stop 获取不到数据时，是否终止当前请求
func (base *BaseController) GetRequestJson(s interface{}, stopRequest bool) error {
	data := base.Ctx.Input.RequestBody
	if len(data) <= 0 || data == nil {
		if stopRequest {
			base.Response(302, "", nil)
			return nil
		} else {
			logs.Info("input is empty")
			return errors.New("input is empty")
		}
	}
	logs.Debug("RequestBody:", utils.BytesToString(data))
	jsonErr := json.Unmarshal(data, s)
	if jsonErr != nil {
		logs.Error("json.Unmarshal is err:", jsonErr.Error())
		base.Response(301, "", nil)
	}
	return nil
}

// 获取登录的用户信息
func (base *BaseController) getLoginUser(uid int64) error {
	userAdministrator, confErr := beego.AppConfig.Int64("admin::userAdministrator")
	if confErr != nil {
		logs.Error(confErr.Error())
		base.Response(500, "", nil)
	}
	var getUserErr error
	utils.CurrentUser.UserModel, getUserErr = models.GetUser(uid) // 查询用户
	if getUserErr != nil {
		logs.Error(getUserErr.Error())
		base.Response(500, "", nil)
		return getUserErr
	}

	if uid == userAdministrator {
		utils.CurrentUser.IsRoot = true
	} else {
		utils.CurrentUser.IsRoot = false
	}
	return nil
}

// 获取limit的默认值
func (base *BaseController) LimitDef(l int) int {
	var Err error
	if l <= 0 {
		limit, _ := beego.AppConfig.String("cms::limit")
		l, Err = strconv.Atoi(limit)
		if Err != nil {
			logs.Error(Err)
			base.Response(500, "", nil)
		}
	}
	return l
}

// 页码默认值
func (base *BaseController) PageDef(p int) int {
	if p <= 0 {
		return 1
	}
	return p
}

// 获取页码偏移量
func (base *BaseController) GetOffset(p int, l int) int {
	return (p - 1) * l
}

// 计算页码偏移量并且检查limit和page
func (base *BaseController) Paginate(p int, l int) (int, int, int) {
	p = base.PageDef(p)
	l = base.LimitDef(l)
	offset := base.GetOffset(p, l)
	return l, p, offset
}

// 初始化权限规则
func (base *BaseController) initRule() {
	// 获取当前的请求URL
	logs.Info(base.Ctx.Input.URI())
	logs.Info(base.Ctx.Input.Param(":id"))
	utils.Check("", utils.CurrentUser.Id, false)
}

//验证表单数据
func (base *BaseController) FormValidation(validData interface{}) {
	v := validation.Validation{}
	b, err := v.Valid(validData)
	if err != nil {
		// handle error
		logs.Error(err.Error())
		return
	}

	//结果验证
	if !b {
		for _, err := range v.Errors {
			logs.Info(err.Key, err.Message)
			base.Response(304, err.Message, nil)
			// 终止app运行
			base.StopRun()
			return
		}
	}
}
