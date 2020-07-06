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
	case 100: //登录成功
		S.Msg = "Login successfully"
		break
	case 101: //登录令牌生成失败
		S.Msg = "Token generation failed"
		break
	case 102: //登录失败，用户名或密码错误
		S.Msg = "Wrong username or password"
		break
	case 103: //身份验证失败，请登录
		S.Msg = "Authentication failed"
		break
	case 104: // 留空
		S.Msg = ""
		break
	case 105: //密码与确认密码不一致
		S.Msg = "New and confirmed passwords are not the same"
		break
	case 106: //登录失败，存在多个相同的用户
		S.Msg = "Login failed, multiple same users exist"
		break
	case 107: //注册失败，用户名或邮箱已被注册
		S.Msg = "User name or mailbox is already registered"
		break
	case 110: //注册成功
		S.Msg = "Registered successfully"
		break
	case 200: //成功码
		S.Msg = "Successful operation"
		break
	case 201: //操作失败
		S.Msg = "Operation failed"
	case 301: //请求数据解析失败
		S.Msg = "Request data parsing failed"
		break
	case 302: //请求数据不能为空
		S.Msg = "Request data cannot be empty"
		break
	case 303: //ID丢失
		S.Msg = "ID missing"
		break
	case 304: //表单验证失败
		S.Msg = "Form validation failed"
		break
	case 400: //数据库查询错误
		S.Msg = "Database query error"
		break
	case 401: //主键丢失
		S.Msg = "Primary key missing"
		break
	case 402: // 存在相同的记录
		S.Msg = "There are the same records"
		break
	case 500: //系统错误
		S.Msg = "System error"
		break
	case 501: //未知的错误
		S.Msg = "Unknown error"
		break
	case 600: //已存在相同的文章栏目别名
		S.Msg = "Duplicate category alias"
		break
	case 601: //找不到栏目
		S.Msg = "Not found category"
		break
	case 602: //找不到文章
		S.Msg = "Not found article"
		break
	case 603: //没有文章被更新
		S.Msg = "Not found article"
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
