package facades

import (
	"go-admin/bean/UserBaen"
	"go-admin/utils"
)

type UserFacade struct {
}

func (receiver *UserFacade) GetUserList() {

}

func (receiver *UserFacade) UserLogin(user UserBaen.LoginJson) error {
	user.Password = utils.Encryption(user.Password)
	return nil
}
