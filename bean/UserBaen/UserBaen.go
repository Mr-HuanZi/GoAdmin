package UserBaen

import "go-admin/models"

var (
	CurrentUser CurrUser
)
//登录表单
type LoginJson struct {
	Username string `valid:"Required;MinSize(4);MaxSize(18);AlphaDash"`
	Password string `valid:"Required;MinSize(6);MaxSize(18);AlphaDash"`
}

//注册表单
type RegisterJson struct {
	Username   string `valid:"Required;MinSize(4);MaxSize(18);AlphaDash"`
	Password   string `valid:"Required;MinSize(6);MaxSize(18);AlphaDash"`
	RePassword string `valid:"Required;MinSize(6);MaxSize(18);AlphaDash"`
	Email      string `valid:"Required;Email"`
}

type CurrUser struct {
	models.UserModel
	IsRoot bool
}
