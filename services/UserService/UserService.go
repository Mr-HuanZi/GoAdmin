package UserService

import (
	"errors"
	"github.com/beego/beego/v2/client/orm"
	"github.com/beego/beego/v2/core/logs"
	beego "github.com/beego/beego/v2/server/web"
	"go-admin/bean/UserBaen"
	"go-admin/models"
	"go-admin/utils"
	"time"
)

var (
	userList = make(map[int64]models.UserModel)
)

// 用户登录
func Login(loginJson UserBaen.LoginJson) (int64, error) {
	loginJson.Password = utils.Encryption(loginJson.Password)
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

func LoginSuccess(uid int64) (string, error) {
	token, tokenErr := utils.GenerateUserToken(uid) //获取登录令牌
	if tokenErr != nil {
		return "", tokenErr //令牌生成失败
	}
	// 获取用户信息
	userAdministrator, confErr := beego.AppConfig.Int64("admin::userAdministrator")
	if confErr != nil {
		return "", confErr
	}
	var getUserErr error
	utils.CurrentUser.UserModel, getUserErr = GetUser(uid) // 查询用户
	if getUserErr != nil {
		return "", getUserErr
	}

	if uid == userAdministrator {
		utils.CurrentUser.IsRoot = true
	} else {
		utils.CurrentUser.IsRoot = false
	}
	//记录用户登录信息
	UpdateUserLoginInfo(uid)
	return token, nil
}

// 获取单个用户信息
func GetUser(uid int64) (models.UserModel, error) {
	_, ok := userList[uid]
	// 从单次缓存中提取用户信息
	if ok {
		return userList[uid], nil
	}
	o := orm.NewOrm()
	user := models.UserModel{Id: uid}
	err := o.Read(&user)
	if err == nil {
		userList[uid] = user // 放入缓存
		return user, nil
	}
	return models.UserModel{}, err
}

func UpdateUserLoginInfo(uid int64) {
	//var ctx *context.Context
	o := orm.NewOrm()
	user := models.UserModel{Id: uid}
	if o.Read(&user) == nil {
		user.LastLoginTime = time.Now().Unix() //获取当前登录时间
		user.LastLoginIp = ""                  //记录登录IP
		user.LockTime = 0
		user.LockTimeStart = 0
		user.ErrorSum = 0
		if num, err := o.Update(&user, "LastLoginTime", "LastLoginIp"); err == nil {
			if num > 0 {
				logs.Info("用户" + user.UserLogin + "登录信息更新成功，登录IP为[" + "127.0.0.1" + "]")
			} else {
				logs.Info("用户" + user.UserLogin + "登录信息更新失败")
			}
		} else {
			logs.Error(err)
		}
	}
}
