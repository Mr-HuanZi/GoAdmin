package admin

import (
	"encoding/json"
	"errors"
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/logs"
	"go-admin/lib/jwt"
	"go-admin/lib/status_code"
	"go-admin/models/admin"
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
