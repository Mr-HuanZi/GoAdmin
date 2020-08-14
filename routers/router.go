// @APIVersion 1.0.0
// @Title beego Test API
// @Description beego has a very cool tools to autogenerate documents for your API
// @Contact astaxie@gmail.com
// @TermsOfServiceUrl http://beego.me/
// @License Apache 2.0
// @LicenseUrl http://www.apache.org/licenses/LICENSE-2.0.html
package routers

import (
	"github.com/astaxie/beego"
	"go-admin/controllers/admin"
	"go-admin/controllers/cms"
)

func init() {
	ns := beego.NewNamespace("/backstage",
		// 登录
		beego.NSRouter("/login", &admin.LoginController{}),
		beego.NSRouter("/login/test", &admin.LoginController{}, "post:Test"),
		beego.NSRouter("/register", &admin.LoginController{}, "post:Register"),
		// 文章
		beego.NSRouter("/article/list", &cms.ArticleController{}, "get,post:List"),
		beego.NSRouter("/article/release", &cms.ArticleController{}, "post:Release"),
		beego.NSRouter("/article/modify", &cms.ArticleController{}, "post:Modify"),
		// 栏目
		beego.NSRouter("/category/list", &cms.CategoryController{}, "get,post:List"),
		beego.NSRouter("/category/add", &cms.CategoryController{}, "post:Add"),
		beego.NSRouter("/category/modify", &cms.CategoryController{}, "post:Modify"),
		// 用户
		beego.NSRouter("/user/list", &admin.UserController{}, "get,post:List"),
		beego.NSRouter("/user/create", &admin.UserController{}, "post:CreateUser"),
		beego.NSRouter("/user/modify", &admin.UserController{}, "post:Modify"),
		// 权限规则
		beego.NSRouter("/rule/list", &admin.RuleController{}, "get,post:List"),
		beego.NSRouter("/rule/add", &admin.RuleController{}, "post:Add"),
		beego.NSRouter("/rule/modify", &admin.RuleController{}, "post:Modify"),
		beego.NSRouter("/rule/WriteGroup", &admin.RuleController{}, "post:WriteGroup"),
		beego.NSRouter("/rule/DeleteGroup", &admin.RuleController{}, "post:DeleteGroup"),
		beego.NSRouter("/rule/AccessAuth", &admin.RuleController{}, "post:AccessAuth"),
		beego.NSRouter("/rule/MemberAuth", &admin.RuleController{}, "post:MemberAuth"),
		beego.NSRouter("/rule/RemoveMemberAuth", &admin.RuleController{}, "post:RemoveMemberAuth"),
	)

	beego.AddNamespace(ns)
}
