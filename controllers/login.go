package controllers

import (
	"go-admin/bean/UserBaen"
	"go-admin/facades"
)

type LoginController struct {
	BaseController
}

//登录
func (c *LoginController) Login() {
	var (
		LoginJson  UserBaen.LoginJson
		UserFacade facades.UserFacade
	)
	_ = c.GetRequestJson(&LoginJson, true)
	/* 表单字段验证 Start */
	c.FormValidation(&LoginJson)
	/* 表单字段验证 End */
	token, loginErr := UserFacade.UserLogin(LoginJson)
	if loginErr != nil {
		c.Response(1, loginErr.Error(), nil)
		return
	}
	c.Ctx.SetCookie("Authorization", token, 7200, "/", "", false, true)
	c.Ctx.Output.Header("Authorization", token)
	c.Response(1, "", nil)
}

//注册用户
func (c *LoginController) Register() {
	var (
		RegisterJson UserBaen.RegisterJson
		UserFacade   facades.UserFacade
	)
	_ = c.GetRequestJson(&RegisterJson, true)
	//验证表单
	c.FormValidation(&RegisterJson)
	_, regErr := UserFacade.Register(RegisterJson)
	if regErr != nil {
		c.Response(0, regErr.Error(), nil)
		return
	}
	c.Response(1, "", nil)
}
