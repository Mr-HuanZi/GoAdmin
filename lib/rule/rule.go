package rule

import (
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/logs"
	"github.com/astaxie/beego/orm"
	"go-admin/models/admin"
)

var (
	AUTH_ON           bool
	AUTH_GROUP        admin.AuthGroupModel       // 权限组表
	AUTH_GROUP_ACCESS admin.AuthGroupAccessModel // 权限组-用户关联表
	AUTH_RULE         admin.RuleModel            // 权限规则表
	AUTH_USER         admin.UserModel            // 用户表
)

func init() {
	var confErr error
	AUTH_ON, confErr = beego.AppConfig.Bool("rule::AUTH_ON")
	if confErr != nil {
		logs.Error(confErr)
		AUTH_ON = false // 如果获取配置失败，则将开关设置为“关”
	}
}

// 验证权限
func Check() {
	logs.Info(AUTH_GROUP)
}

// 根据用户id获取用户组,返回值为数组
func getGroup(uid int64) []int {
	var (
		groupData []int
	)
	// 判断用户是否存在
	user := admin.UserModel{Id: uid}
	o := orm.NewOrm()
	userReadErr := o.Read(&user)

	if userReadErr == orm.ErrNoRows {
		logs.Error("权限组查不到该用户")
	} else if userReadErr == orm.ErrMissPK {
		logs.Error("找不到主键")
	} else {
		// 获取关联表的数据
		var (
			accessData []*admin.AuthGroupAccessModel
			accessErr  error
		)
		_, accessErr = o.QueryTable(AUTH_GROUP_ACCESS).Filter("uid", user.Id).All(&accessData)
		if accessErr != nil {
			logs.Error(accessErr)
		}
		groupData = append(groupData, len(accessData))
		for k, v := range accessData {
			groupData[k] = v.GroupId
		}
	}
	return groupData
}

// 根据用户ID获取用户规则列表，返回值为数组
func getUserRuleList() {
}
