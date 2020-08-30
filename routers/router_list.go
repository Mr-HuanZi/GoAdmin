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
	{
		Router:     "/article/list",
		Controller: &cms.ArticleController{},
		Method:     "get,post:List",
	},
	{
		Router:     "/article/release",
		Controller: &cms.ArticleController{},
		Method:     "post:Release",
	},
	{
		Router:     "/article/modify",
		Controller: &cms.ArticleController{},
		Method:     "post:Modify",
	},
	{
		Router:     "/article/delete",
		Controller: &cms.ArticleController{},
		Method:     "post:Delete",
	},
	//************************* 栏目 *************************//
	{
		Router:     "/category/list",
		Controller: &cms.CategoryController{},
		Method:     "get,post:List",
	},
	{
		Router:     "/category/add",
		Controller: &cms.CategoryController{},
		Method:     "post:Add",
	},
	{
		Router:     "/category/modify",
		Controller: &cms.CategoryController{},
		Method:     "post:Modify",
	},
	{
		Router:     "/category/delete",
		Controller: &cms.CategoryController{},
		Method:     "post:Delete",
	},
	//************************* 用户 *************************//
	{
		Router:     "/user/list",
		Controller: &admin.UserController{},
		Method:     "get,post:List",
	},
	{
		Router:     "/user/create",
		Controller: &admin.UserController{},
		Method:     "post:CreateUser",
	},
	{
		Router:     "/user/modify",
		Controller: &admin.UserController{},
		Method:     "post:Modify",
	},
	{
		Router:     "/user/forbid",
		Controller: &admin.UserController{},
		Method:     "post:ForbidUser",
	},
	{
		Router:     "/user/resume",
		Controller: &admin.UserController{},
		Method:     "post:ResumeUser",
	},
	//************************* 权限规则 *************************//
	{
		Router:     "/rule/list",
		Controller: &admin.RuleController{},
		Method:     "get,post:List",
	},
	{
		Router:     "/rule/add",
		Controller: &admin.RuleController{},
		Method:     "post:Add",
	},
	{
		Router:     "/rule/modify",
		Controller: &admin.RuleController{},
		Method:     "post:Modify",
	},
	{
		Router:     "/rule/WriteGroup",
		Controller: &admin.RuleController{},
		Method:     "post:WriteGroup",
	},
	{
		Router:     "/rule/DeleteGroup",
		Controller: &admin.RuleController{},
		Method:     "post:DeleteGroup",
	},
	{
		Router:     "/rule/AccessAuth",
		Controller: &admin.RuleController{},
		Method:     "post:AccessAuth",
	},
	{
		Router:     "/rule/MemberAuth",
		Controller: &admin.RuleController{},
		Method:     "post:MemberAuth",
	},
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
