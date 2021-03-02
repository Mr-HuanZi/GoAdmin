package sys_logs

import (
	"encoding/json"
	"github.com/beego/beego/v2/client/orm"
	"github.com/beego/beego/v2/core/logs"
	"github.com/beego/beego/v2/server/web/context"
	"go-admin/lib"
	"go-admin/models/admin"
	"time"
)

// 写入日志操作
// @logType 日志类型 1-登录日志 2-操作日志
// @content 日志详情
func WriteSysLogs(logType int8, logText string, Input *context.BeegoInput, more string) {
	logsData := admin.SystemLogsModel{
		Type:       logType,
		CreateTime: time.Now().Unix(),
		Ip:         Input.IP(),
		Url:        Input.URI(),
		UserId:     0,
		Username:   "",
		Content:    logText,
		Param:      "",
	}

	if lib.CurrentUser.Id > 0 {
		logsData.UserId = lib.CurrentUser.Id
		logsData.Username = lib.CurrentUser.UserLogin
	}

	// 处理请求数据
	param := make(map[string]interface{})
	param["requestQuery"] = ""

	// 处理请求Body
	if len(Input.RequestBody) > 0 {
		requestBody := make(map[string]interface{})
		// 请求的body
		jsonDnErr := json.Unmarshal(Input.RequestBody, &requestBody)
		if jsonDnErr != nil {
			logs.Warning("json.Unmarshal is err:", jsonDnErr.Error())
		}
		param["requestBody"] = requestBody
	}

	// 解析more
	if more != "" {
		moreMap := make(map[string]interface{})
		jsonDnErr := json.Unmarshal([]byte(more), &moreMap)
		if jsonDnErr != nil {
			logs.Warning("json.Unmarshal is err:", jsonDnErr.Error())
			param["more"] = ""
		} else {
			param["more"] = moreMap
		}
	}

	moreJson, jsonMarshalErr := json.Marshal(param)
	if jsonMarshalErr != nil {
		logs.Warning("json.Marshal is err:", jsonMarshalErr.Error())
		logsData.Param = ""
	} else {
		logsData.Param = lib.BytesToString(moreJson)
	}

	// 写库
	o := orm.NewOrm()
	_, err := o.Insert(&logsData)
	if err != nil {
		logs.Error(err)
	}
}
