package routers

import (
	"github.com/astaxie/beego"
	"go-admin/controllers/admin"
	"go-admin/controllers/cms"
)

type RouterList struct {
	Router     string
	Controller beego.ControllerInterface
	Method     string
}

var RouterListInterface = []RouterList{
	//************************* 登录 *************************//
	{
		Router:     "/login",
		Controller: &admin.LoginController{},
		Method:     "post:Login",
	},
	{
		Router:     "/login/test",
		Controller: &admin.LoginController{},
		Method:     "post:Test",
	},
	{
		Router:     "/register",
		Controller: &admin.LoginController{},
		Method:     "post:Register",
	},
	//************************* 文章 *************************//
	// 文章列表
	{
		Router:     "/article/list",
		Controller: &cms.ArticleController{},
		Method:     "get,post:List",
	},
	// 文章发布
	{
		Router:     "/article/release",
		Controller: &cms.ArticleController{},
		Method:     "post:Release",
	},
	// 文章修改
	{
		Router:     "/article/modify",
		Controller: &cms.ArticleController{},
		Method:     "post:Modify",
	},
	// 文章删除
	{
		Router:     "/article/delete",
		Controller: &cms.ArticleController{},
		Method:     "post:Delete",
	},
	//************************* 栏目 *************************//
	// 栏目列表
	{
		Router:     "/category/list",
		Controller: &cms.CategoryController{},
		Method:     "get,post:List",
	},
	// 新增栏目
	{
		Router:     "/category/add",
		Controller: &cms.CategoryController{},
		Method:     "post:Add",
	},
	// 栏目修改
	{
		Router:     "/category/modify",
		Controller: &cms.CategoryController{},
		Method:     "post:Modify",
	},
	// 栏目删除
	{
		Router:     "/category/delete",
		Controller: &cms.CategoryController{},
		Method:     "post:Delete",
	},
	//************************* 用户 *************************//
	// 用户列表
	{
		Router:     "/user/list",
		Controller: &admin.UserController{},
		Method:     "get,post:List",
	},
	// 创建用户
	{
		Router:     "/user/create",
		Controller: &admin.UserController{},
		Method:     "post:CreateUser",
	},
	// 用户信息修改
	{
		Router:     "/user/modify",
		Controller: &admin.UserController{},
		Method:     "post:Modify",
	},
	// 禁用用户
	{
		Router:     "/user/forbid",
		Controller: &admin.UserController{},
		Method:     "post:ForbidUser",
	},
	// 启用用户
	{
		Router:     "/user/resume",
		Controller: &admin.UserController{},
		Method:     "post:ResumeUser",
	},
	{
		Router:     "/user/fetchCurrent",
		Controller: &admin.UserController{},
		Method:     "get:FetchCurrentUser",
	},
	//************************* 权限规则 *************************//
	// 权限规则列表
	{
		Router:     "/rule/list",
		Controller: &admin.RuleController{},
		Method:     "get,post:List",
	},
	// 新增规则
	{
		Router:     "/rule/add",
		Controller: &admin.RuleController{},
		Method:     "post:Add",
	},
	// 修改规则
	{
		Router:     "/rule/modify",
		Controller: &admin.RuleController{},
		Method:     "post:Modify",
	},
	// 写入权限组
	{
		Router:     "/rule/WriteGroup",
		Controller: &admin.RuleController{},
		Method:     "post:WriteGroup",
	},
	// 删除权限组
	{
		Router:     "/rule/DeleteGroup",
		Controller: &admin.RuleController{},
		Method:     "post:DeleteGroup",
	},
	// 权限组授权
	{
		Router:     "/rule/AccessAuth",
		Controller: &admin.RuleController{},
		Method:     "post:AccessAuth",
	},
	// 用户授权
	{
		Router:     "/rule/MemberAuth",
		Controller: &admin.RuleController{},
		Method:     "post:MemberAuth",
	},
	// 移除人员权限
	{
		Router:     "/rule/RemoveMemberAuth",
		Controller: &admin.RuleController{},
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
