package facades

import (
	"github.com/beego/beego/v2/core/logs"
	"go-admin/bean/UserBaen"
	"go-admin/services/UserService"
)

type UserFacade struct {
}

func (receiver *UserFacade) GetUserList() {

}

func (receiver *UserFacade) UserLogin(user UserBaen.LoginJson) error {
	uid, err := UserService.Login(user)
	if err != nil {
		return err
	}
	logs.Info(uid)
	return nil
}
