package routers

import (
	beego "github.com/beego/beego/v2/server/web"
	"go-admin/controllers"
)

type RouterList struct {
	Router     string
	Controller beego.ControllerInterface
	Method     string
}

var RouterListInterface = []RouterList{
	//************************* 公共 *************************//
	{
		// 单文件上传
		Router:     "/upload",
		Controller: &controllers.FileUploadController{},
		Method:     "post:UploadFile",
	},
	//************************* 登录 *************************//
	{
		Router:     "/login",
		Controller: &controllers.LoginController{},
		Method:     "post:Login",
	},
	{
		Router:     "/register",
		Controller: &controllers.LoginController{},
		Method:     "post:Register",
	},
	//************************* 文章 *************************//
	// 文章列表
	{
		Router:     "/article/list",
		Controller: &controllers.ArticleController{},
		Method:     "get:List",
	},
	// 文章发布
	{
		Router:     "/article/release",
		Controller: &controllers.ArticleController{},
		Method:     "post:Release",
	},
	// 文章修改
	{
		Router:     "/article/modify",
		Controller: &controllers.ArticleController{},
		Method:     "post:Modify",
	},
	// 文章删除
	{
		Router:     "/article/delete",
		Controller: &controllers.ArticleController{},
		Method:     "get:Delete",
	},
	// 获取单篇文章
	{
		Router:     "/article/one",
		Controller: &controllers.ArticleController{},
		Method:     "get:GetArticle",
	},
	//************************* 栏目 *************************//
	// 栏目列表
	{
		Router:     "/category/list",
		Controller: &controllers.CategoryController{},
		Method:     "get,post:List",
	},
	// 新增栏目
	{
		Router:     "/category/create",
		Controller: &controllers.CategoryController{},
		Method:     "post:Add",
	},
	// 栏目修改
	{
		Router:     "/category/modify",
		Controller: &controllers.CategoryController{},
		Method:     "post:Modify",
	},
	// 栏目删除
	{
		Router:     "/category/delete",
		Controller: &controllers.CategoryController{},
		Method:     "post:Delete",
	},
	// 获取单个栏目
	{
		Router:     "/category/:id:int",
		Controller: &controllers.CategoryController{},
		Method:     "get:GetCategory",
	},
	//************************* 用户 *************************//
	// 用户列表
	{
		Router:     "/user/list",
		Controller: &controllers.UserController{},
		Method:     "get,post:List",
	},
	// 创建用户
	{
		Router:     "/user/create",
		Controller: &controllers.UserController{},
		Method:     "post:CreateUser",
	},
	// 用户信息修改
	{
		Router:     "/user/modify",
		Controller: &controllers.UserController{},
		Method:     "post:Modify",
	},
	// 禁用用户
	{
		Router:     "/user/forbid",
		Controller: &controllers.UserController{},
		Method:     "post:ForbidUser",
	},
	// 启用用户
	{
		Router:     "/user/resume",
		Controller: &controllers.UserController{},
		Method:     "post:ResumeUser",
	},
	{
		Router:     "/user/fetchCurrent",
		Controller: &controllers.UserController{},
		Method:     "get:FetchCurrentUser",
	},
	//************************* 权限规则 *************************//
	// 权限规则列表
	{
		Router:     "/rule/list",
		Controller: &controllers.RuleController{},
		Method:     "get,post:List",
	},
	// 新增规则
	{
		Router:     "/rule/add",
		Controller: &controllers.RuleController{},
		Method:     "post:Add",
	},
	// 修改规则
	{
		Router:     "/rule/modify",
		Controller: &controllers.RuleController{},
		Method:     "post:Modify",
	},
	// 写入权限组
	{
		Router:     "/rule/WriteGroup",
		Controller: &controllers.RuleController{},
		Method:     "post:WriteGroup",
	},
	// 删除权限组
	{
		Router:     "/rule/DeleteGroup",
		Controller: &controllers.RuleController{},
		Method:     "post:DeleteGroup",
	},
	// 权限组授权
	{
		Router:     "/rule/AccessAuth",
		Controller: &controllers.RuleController{},
		Method:     "post:AccessAuth",
	},
	// 用户授权
	{
		Router:     "/rule/MemberAuth",
		Controller: &controllers.RuleController{},
		Method:     "post:MemberAuth",
	},
	// 移除人员权限
	{
		Router:     "/rule/RemoveMemberAuth",
		Controller: &controllers.RuleController{},
		Method:     "post:RemoveMemberAuth",
	},
}

func GetRouterList() []beego.LinkNamespace {
	var l []beego.LinkNamespace
	for _, list := range RouterListInterface {
		l = append(l, beego.NSRouter(list.Router, list.Controller, list.Method))
	}
	return l
}
