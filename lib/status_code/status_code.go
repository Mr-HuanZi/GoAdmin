package status_code

type StatusCode struct {
	Code int
	Msg  string
	Data interface{}
}

// 根据状态码获取文本
// 代码段说明
// 1 登录/注册
// 2 操作结果类
// 3 请求处理类
// 4 数据库类
// 5 系统响应类
// 6 app响应类
func (S *StatusCode) GetStatusCode(code int) *StatusCode {
	S.Code = code
	switch code {
	// 登录代码
	case 100: //登录成功
		S.Msg = "登录成功"
		break
	case 101: //登录令牌生成失败
		S.Msg = "登录令牌生成失败"
		break
	case 102: //登录失败，用户名或密码错误
		S.Msg = "用户名或密码错误"
		break
	case 103: //身份验证失败，请登录
		S.Msg = "无效的令牌"
		break
	case 104: // 用户未启用，禁止登录
		S.Msg = "用户未启用"
		break
	case 105: //密码与确认密码不一致
		S.Msg = "新密码与确认密码不一致"
		break
	case 106: //登录失败，存在多个相同的用户
		S.Msg = "登录失败，存在多个相同的用户"
		break
	case 107: //注册失败，用户名或邮箱已被注册
		S.Msg = "用户名或邮箱已被注册"
		break
	case 108: //登录失败，用户被锁定
		S.Msg = "账户被锁定"
		break
	case 110: //注册成功
		S.Msg = "注册成功"
		break
	// 操作类代码
	case 200: //成功码
		S.Msg = "操作成功"
		break
	case 201: //操作失败
		S.Msg = "操作失败"
	// 请求类代码
	case 301: //请求数据解析失败
		S.Msg = "请求数据解析失败"
		break
	case 302: //请求数据不能为空
		S.Msg = "请求数据不能为空"
		break
	case 303: //ID丢失
		S.Msg = "ID丢失"
		break
	case 304: //表单验证失败
		S.Msg = "表单验证失败"
		break
	// 数据库类代码
	case 400: //数据库查询错误
		S.Msg = "数据库查询错误"
		break
	case 401: //主键丢失
		S.Msg = "主键丢失"
		break
	case 402: // 存在相同的记录
		S.Msg = "存在相同的记录"
		break
	case 403: // 没有记录被更新
		S.Msg = "没有记录被更新"
		break
	case 404: // 查询不到结果
		S.Msg = "查询不到结果"
		break
	case 405: // 没有记录被删除
		S.Msg = "没有记录被删除"
		break
	// 系统级别代码
	case 500: // 系统错误
		S.Msg = "系统错误"
		break
	case 501: // 未知的错误
		S.Msg = "未知错误"
		break
	case 502: // 数据格式不正确
		S.Msg = "数据格式不正确"
		break
	case 503: // 数据类型错误
		S.Msg = "数据类型错误"
		break
	// 模块类代码
	case 600: //已存在相同的文章栏目别名
		S.Msg = "已存在相同的文章栏目别名"
		break
	case 601: //找不到栏目
		S.Msg = "找不到栏目"
		break
	case 602: //找不到文章
		S.Msg = "找不到文章"
		break
	}
	return S
}

func (S *StatusCode) CreateData(code int, msg string, data interface{}) *StatusCode {
	S.GetStatusCode(code)
	if msg != "" {
		S.Msg = msg
	}
	S.Data = &data
	return S
}
