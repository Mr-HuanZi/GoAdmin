package admin

import (
	"encoding/json"
	"errors"
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/logs"
	"go-admin/lib/jwt"
	"go-admin/lib/rule"
	"go-admin/lib/status_code"
	"go-admin/models/admin"
	"strconv"
	"strings"
)

type BaseController struct {
	beego.Controller
	// 当前登录的用户信息
	ThatUser admin.UserModel
	UserRoot bool // 当前登录用户是否管理员
}

//实例化日志模块
var StatusCodeInstance status_code.StatusCode

// 基类初始化方法，每次执行都会被调用
func (base *BaseController) Prepare() {
	// 日志开始
	logs.Info("\r\n\r\n")
	logs.Info(base.Ctx.Input.IP(), "["+base.Ctx.Input.Method()+"]", base.Ctx.Input.URL(), base.Ctx.Input.URI())
	//验证登录
	var checkLogin = true
	var ruleExclude []string
	ruleExcludeConf := beego.AppConfig.String("rule::ruleExclude")
	if i := strings.Index(ruleExcludeConf, ","); i != -1 {
		ruleExclude = strings.Split(ruleExcludeConf, ",")
	} else {
		ruleExclude = append(ruleExclude, ruleExcludeConf)
	}
	for _, value := range ruleExclude {
		if value == base.Ctx.Request.RequestURI {
			checkLogin = false
			break
		}
	}
	if checkLogin {
		Authorization := base.Ctx.GetCookie("Authorization")
		b, t := jwt.ValidateUserToken(Authorization)
		if !b {
			//token验证失败
			base.Response(103, "", nil)
		} else {
			// 获取当前登录用户信息
			LoginToken := jwt.GetTokenClaims(t)
			getUserErr := base.getLoginUser(LoginToken.User)
			if getUserErr != nil {
				base.Response(500, getUserErr.Error(), nil) //令牌生成失败
			}
			// 刷新令牌
			token, tokenErr := jwt.RefreshUserToken(t)
			if tokenErr != nil {
				base.Response(101, "", nil) //令牌生成失败
			}
			if token != "" {
				base.Ctx.SetCookie("Authorization", token, 3600, "/", "", false, true)
			}
		}
	}

	// 初始化权限规则
	initRule()
}

func (base *BaseController) Response(code int, msg string, data interface{}) {
	statusCode := StatusCodeInstance.CreateData(code, msg, data)
	logs.Info(statusCode)
	base.Data["json"] = &statusCode
	base.ServeJSON()
	base.StopRun()
}

// 获取json请求数据
// @param s 接收数据的结构体
// @param stop 获取不到数据时，是否终止当前请求
func (base *BaseController) GetRequestJson(s interface{}, stopRequest bool) error {
	data := base.Ctx.Input.RequestBody
	if len(data) <= 0 || data == nil {
		if stopRequest {
			base.Response(302, "", nil)
		} else {
			logs.Info("input is empty")
			return errors.New("input is empty")
		}
	}
	logs.Info("RequestBody:")
	logs.Info(data)
	jsonErr := json.Unmarshal(data, s)
	if jsonErr != nil {
		logs.Error("json.Unmarshal is err:", jsonErr.Error())
		base.Response(301, "", nil)
	}
	logs.Info(s)
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
	base.ThatUser, getUserErr = admin.GetUser(uid) // 查询用户
	if getUserErr != nil {
		logs.Error(getUserErr.Error())
		base.Response(500, "", nil)
		return getUserErr
	}

	if uid == userAdministrator {
		base.UserRoot = true
	} else {
		base.UserRoot = false
	}
	return nil
}

// 获取limit的默认值
func (base *BaseController) LimitDef(l int) int {
	var Err error
	if l <= 0 {
		limit := beego.AppConfig.String("cms::limit")
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
func initRule() {
	rule.Check()
}
