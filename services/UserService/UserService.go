package UserService

import (
	"errors"
	"github.com/beego/beego/v2/client/orm"
	"go-admin/bean/UserBaen"
	"go-admin/models"
	"time"
)

// 用户登录
func login(loginJson UserBaen.LoginJson) (int64, error) {
	o := orm.NewOrm()
	var user models.UserModel
	//查询完整的记录
	err := o.QueryTable(user).Filter("user_type", 1).Filter("user_login", loginJson.Username).One(&user)
	if err == orm.ErrMultiRows {
		// 查询到多个记录
		return 0, errors.New("查询到多个记录")
	} else if user.Id != 0 {
		// 确认用户是否被锁定
		if user.LockTime > 0 && (time.Now().Unix() < user.LockTime) {
			return 0, errors.New("用户已被锁定")
		}
		// 检查密码是否正确
		if loginJson.Password != user.UserPass {
			// 记录错误次数
			models.AddUserLoginErrorSum(&user)
			if user.ErrorSum+1 >= 3 {
				return 0, errors.New("用户已被锁定")
			}
			return 0, errors.New("账号或密码错误")
		} else if user.UserStatus != 1 {
			// 账户是否是“启用”状态
			return 0, errors.New("用户未启用")
		}
		return user.Id, nil
	}
	return 0, errors.New("账号或密码错误")
}
