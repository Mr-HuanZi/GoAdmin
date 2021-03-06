package utils

import (
	"errors"
	"github.com/beego/beego/v2/client/orm"
	"github.com/beego/beego/v2/core/logs"
	beego "github.com/beego/beego/v2/server/web"
	"go-admin/models"
	"strings"
)

var (
	AuthON          bool
	AuthGroup       models.AuthGroupModel       // 权限组表
	AuthGroupAccess models.AuthGroupAccessModel // 权限组-用户关联表
	AuthRule        models.RuleModel            // 权限规则表
	AuthUser        models.UserModel            // 用户表
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
func Check(rule string, uid int64, strict bool) bool {
	if !AuthON {
		return true
	}
	// 获取当前用户的规则列表
	ruleList, _ := getUserRuleList(uid)
	logs.Info(ruleList)
	//
	return true
}

// 根据用户id获取用户组,返回值为数组
func getGroup(uid int64, groupData interface{}) (int64, error) {
	var (
		queryNum int64
	)
	// 判断用户是否存在
	user := models.UserModel{Id: uid}
	o := orm.NewOrm()
	userReadErr := o.Read(&user)

	if userReadErr == orm.ErrNoRows {
		logs.Error("权限组查不到该用户")
	} else if userReadErr == orm.ErrMissPK {
		logs.Error("找不到主键")
	} else {
		// 获取关联表的数据
		var (
			accessData []*models.AuthGroupAccessModel
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
	var (
		groupData []*models.AuthGroupModel
		ruleList  []*models.RuleModel
		ruleStr   []string
	)
	groupNum, groupErr := getGroup(uid, &groupData)
	if groupErr != nil {
		return nil, groupErr
	}
	if groupNum > 0 {
		var ruleIds []string
		for _, group := range groupData {
			logs.Info(group.Rules)
			rule := strings.Split(strings.TrimSpace(group.Rules), ",")
			ruleIds = append(ruleIds, rule...)
		}
		// 如果获取数据失败
		if len(ruleIds) <= 0 {
			return nil, errors.New("ruleIds failed to get data")
		}
		// 从规则表中获取数据
		o := orm.NewOrm()
		ruleModelQueryNum, ruleModelQueryErr := o.QueryTable(AuthRule).Filter("id__in", ruleIds).All(&ruleList)
		if ruleModelQueryErr != nil {
			return nil, ruleModelQueryErr
		}
		if ruleModelQueryNum > 0 {
			for _, rule := range ruleList {
				var s string
				if rule.Param != "" {
					if strings.Index(rule.Param, "?") != -1 {
						s = rule.Rule + rule.Param
					} else {
						s = rule.Rule + "?" + rule.Param
					}
				} else {
					s = rule.Rule
				}
				ruleStr = append(ruleStr, s)
			}
			return ruleStr, nil
		}
		return nil, errors.New("not found")
	}
	return nil, nil
}
