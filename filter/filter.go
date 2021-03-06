package filter

import (
	"github.com/beego/beego/v2/core/logs"
	"github.com/beego/beego/v2/server/web"
	"github.com/beego/beego/v2/server/web/context"
	"go-admin/utils"
)

func ExecAdminFilter() {
	// 静态地址之前执行
	web.InsertFilter("*", web.BeforeStatic, func(ctx *context.Context) {
		// 日志开始
		logs.Info("------------------------------ Controller run ["+ctx.Input.Method()+"] " + ctx.Input.IP() + ctx.Input.URI() +" ------------------------------")
	})
	// 执行完 Controller 逻辑之后执行的过滤器
	web.InsertFilter("*", web.AfterExec, func(ctx *context.Context) {
		// 把当前请求数据写入到数据库
		go utils.WriteSysLogs(1, "RequestInput", ctx.Input, "")
	})
}
