package rule

import (
	"errors"
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/logs"
	"github.com/astaxie/beego/orm"
	"go-admin/models/admin"
	"strings"
)

var (
	AuthON          bool
	AuthGroup       admin.AuthGroupModel       // 权限组表
	AuthGroupAccess admin.AuthGroupAccessModel // 权限组-用户关联表
	AuthRule        admin.RuleModel            // 权限规则表
	AuthUser        admin.UserModel            // 用户表
)

func init() {
	var confErr error
	AuthON, confErr = beego.AppConfig.Bool("rule::AUTH_ON")
	if confErr != nil {
		logs.Error(confErr)
		AuthON = false // 如果获取配置失败，则将开关设置为“关”
	}
}

// 验证权限
func Check() bool {
	if !AuthON {
		return true
	}
_:
	getUserRuleList(1)
	return true
}

// 根据用户id获取用户组,返回值为数组
func getGroup(uid int64, groupData interface{}) (int64, error) {
	var (
		queryNum int64
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
			queryErr   error
		)
		queryNum, queryErr = o.QueryTable(AuthGroupAccess).Filter("uid", user.Id).All(&accessData)
		if queryErr != nil {
			logs.Error(queryErr)
			return 0, queryErr
		}
		// 如果返回的记录数大于0
		if queryNum > 0 {
			logs.Info(queryNum)
			// 获取组ID
			groupIds := make([]int, len(accessData))
			for k, v := range accessData {
				groupIds[k] = v.GroupId
			}
			logs.Info(groupIds)
			// 获取每个组的数据
			queryNum, queryErr = o.QueryTable(AuthGroup).Filter("id__in", groupIds).All(groupData)
			if queryErr != nil {
				logs.Error(queryErr)
				return 0, queryErr
			}
		} else {
			return 0, errors.New("no data")
		}
	}
	return queryNum, nil
}

// 根据用户ID获取用户规则列表，返回值为数组
func getUserRuleList(uid int64) ([]string, error) {
	var groupData []*admin.AuthGroupModel
	groupNum, groupErr := getGroup(uid, &groupData)
	if groupErr != nil {
		return nil, groupErr
	}
	if groupNum > 0 {
		var ruleIds []string
		for _, group := range groupData {
			rule := strings.Split(strings.TrimSpace(group.Rules), ",")
			ruleIds = append(ruleIds, rule...)
		}
		// 如果获取数据失败
		if len(ruleIds) <= 0 {
			return nil, errors.New("ruleIds failed to get data")
		}
	}

	return nil, nil
}
