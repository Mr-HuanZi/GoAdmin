package utils

import (
	"github.com/beego/beego/v2/core/logs"
	"github.com/beego/beego/v2/server/web"
	"strings"
)

//  检查rule是否在配置ruleExclude中
func InRuleExclude(rule string) bool {
	var ruleExclude []string
	ruleExcludeConf, err := web.AppConfig.String("rule::ruleExclude")
	if err != nil {
		logs.Error(err.Error())
		return false
	}
	if ruleExcludeConf == "" {
		return true
	}
	if i := strings.Index(ruleExcludeConf, ","); i != -1 {
		ruleExclude = strings.Split(ruleExcludeConf, ",")
	} else {
		ruleExclude = append(ruleExclude, ruleExcludeConf)
	}
	for _, value := range ruleExclude {
		if strings.ToUpper(value) == strings.ToUpper(rule) {
			return true
		}
	}
	return false
}
