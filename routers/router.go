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
		beego.NSRouter("/login", &admin.LoginController{}),
		beego.NSRouter("/login/test", &admin.LoginController{}, "post:Test"),
		beego.NSRouter("/register", &admin.LoginController{}, "post:Register"),
		beego.NSRouter("/article/list", &cms.ArticleController{}, "post:List"),
		beego.NSRouter("/article/release", &cms.ArticleController{}, "post:Release"),
		beego.NSRouter("/article/modify", &cms.ArticleController{}, "post:Modify"),
		beego.NSRouter("/category/list", &cms.CategoryController{}, "post:List"),
		beego.NSRouter("/category/add", &cms.CategoryController{}, "post:Add"),
		beego.NSRouter("/category/modify", &cms.CategoryController{}, "post:Modify"),
		beego.NSRouter("/user/list", &admin.UserController{}, "post:List"),
		beego.NSRouter("/user/create", &admin.UserController{}, "post:CreateUser"),
	)

	beego.AddNamespace(ns)
}
