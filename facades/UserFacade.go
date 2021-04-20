package facades

import (
	"go-admin/bean/UserBaen"
	"go-admin/services"
)

type UserFacade struct {
}

func (receiver *UserFacade) GetUserList() {
}

// 用户登录
func (receiver *UserFacade) UserLogin(user UserBaen.LoginJson) (string, error) {
	UserService := services.GetUserServiceInstance()
	uid, err := UserService.Login(user)
	if err != nil {
		return "", err
	}
	return UserService.LoginSuccess(uid)
}

// 注册用户
func (receiver UserFacade) Register(regData UserBaen.RegisterJson) (int64, error) {
	UserService := services.GetUserServiceInstance()
	return UserService.Register(regData)
}
