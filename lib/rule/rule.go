package rule

import (
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/logs"
)

var (
	AUTH_ON           bool
	AUTH_GROUP        string
	AUTH_GROUP_ACCESS string
	AUTH_RULE         string
	AUTH_USER         string
)

func init() {
	var confErr error
	AUTH_ON, confErr = beego.AppConfig.Bool("rule::AUTH_ON")
	if confErr != nil {
		logs.Error(confErr)
	}
	AUTH_GROUP = beego.AppConfig.String("rule::AUTH_GROUP")
	AUTH_GROUP_ACCESS = beego.AppConfig.String("rule::AUTH_GROUP_ACCESS")
	AUTH_RULE = beego.AppConfig.String("rule::AUTH_RULE")
	AUTH_USER = beego.AppConfig.String("rule::AUTH_USER")
}

// 验证权限
func Check() {
}

// 根据用户id获取用户组,返回值为数组
func getGroup(uid int64) {
}

// 根据用户ID获取用户规则列表，返回值为数组
func getUserRuleList() {
}
