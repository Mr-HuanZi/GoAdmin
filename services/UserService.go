package services

import (
	"errors"
	"github.com/beego/beego/v2/client/orm"
	"github.com/beego/beego/v2/core/logs"
	beego "github.com/beego/beego/v2/server/web"
	"go-admin/bean/Request"
	"go-admin/bean/UserBaen"
	"go-admin/models"
	"go-admin/utils"
	"sync"
	"time"
)

var (
	userList = make(map[int64]models.UserModel)
	instance *UserService
	once     sync.Once
)

type UserService struct {
}

// 单例模式
func GetUserServiceInstance() *UserService {
	once.Do(func() {
		instance = &UserService{}
	})
	return instance
}

// 用户登录
func (that *UserService) Login(loginJson UserBaen.LoginJson) (int64, error) {
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
			that.AddUserLoginErrorSum(&user)
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

func (that *UserService) LoginSuccess(uid int64) (string, error) {
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
	UserBaen.CurrentUser.UserModel, getUserErr = that.GetUser(uid) // 查询用户
	if getUserErr != nil {
		return "", getUserErr
	}

	if uid == userAdministrator {
		UserBaen.CurrentUser.IsRoot = true
	} else {
		UserBaen.CurrentUser.IsRoot = false
	}
	//记录用户登录信息
	that.UpdateUserLoginInfo(uid)
	return token, nil
}

// 获取单个用户信息
func (that *UserService) GetUser(uid int64) (models.UserModel, error) {
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

func (that *UserService) UpdateUserLoginInfo(uid int64) {
	o := orm.NewOrm()
	ip := Request.InputCopy.IP()
	user := models.UserModel{Id: uid}
	if o.Read(&user) == nil {
		user.LastLoginTime = time.Now().Unix() //获取当前登录时间
		user.LastLoginIp = ip                  //记录登录IP
		user.LockTime = 0
		user.LockTimeStart = 0
		user.ErrorSum = 0
		if num, err := o.Update(&user, "LastLoginTime", "LastLoginIp"); err == nil {
			if num > 0 {
				logs.Info("用户" + user.UserLogin + "登录信息更新成功，登录IP为[" + ip + "]")
			} else {
				logs.Info("用户" + user.UserLogin + "登录信息更新失败")
			}
		} else {
			logs.Error(err)
		}
	}
}

// 注册用户
func (that *UserService) Register(registerJson UserBaen.RegisterJson) (int64, error) {
	//确认密码是否相等
	if registerJson.Password != registerJson.RePassword {
		return 0, errors.New("设置的密码与确认密码不一致")
	}
	//检查用户名和邮箱是否已被注册
	if that.CheckUserRepeat(registerJson.Username, registerJson.Email) {
		return 0, errors.New("用户名或邮箱已被注册")
	}
	//创建用户数据
	var UserData models.UserModel
	var uid int64
	UserData.UserLogin = registerJson.Username
	UserData.UserPass = utils.Encryption(registerJson.Password) //加密密码
	UserData.UserNickname = registerJson.Username               //用户昵称默认是登录账号
	UserData.UserType = 1                                       //管理员类型
	UserData.CreateTime = time.Now().Unix()                     //管理员类型
	UserData.UpdateTime = UserData.CreateTime                   //管理员类型
	UserData.UserStatus = 1                                     //用户状态
	uid = that.CreateUser(&UserData)
	if uid != 0 {
		return uid, nil
	}
	return 0, errors.New("注册失败")
}

// 记录用户登录错误次数
func (that *UserService) AddUserLoginErrorSum(user *models.UserModel) {
	if user.Id <= 0 {
		return
	}
	var (
		err error
		num int64
	)
	o := orm.NewOrm()
	if user.ErrorSum <= 3 {
		// 登录错误次数小于3次，只记录错误次数
		num, err = o.QueryTable(new(models.UserModel)).Filter("id", user.Id).Update(orm.Params{
			"ErrorSum": orm.ColValue(orm.ColAdd, 1),
		})
	} else {
		// 获取当前时间
		nowTime := time.Now()
		// 时间增加1小时
		h, _ := time.ParseDuration("1h")
		lockTime := nowTime.Add(h).Unix()
		num, err = o.QueryTable(new(models.UserModel)).Filter("id", user.Id).Update(orm.Params{
			"error_sum":       orm.ColValue(orm.ColAdd, 1),
			"lock_time":       lockTime,
			"lock_time_start": nowTime.Unix(),
		})
	}
	if err != nil {
		logs.Error(err)
	}
	logs.Debug("AddUserLoginErrorSum:", num)
}

//创建新用户
func (that *UserService) CreateUser(userData *models.UserModel) int64 {
	o := orm.NewOrm()
	id, err := o.Insert(userData)
	if err == nil {
		logs.Info(id)
		return id
	}
	return 0
}

// 检查用户是否重复
func (that *UserService) CheckUserRepeat(username string, email string) bool {
	o := orm.NewOrm()
	//自定义条件
	ormCondition := orm.NewCondition()
	ormConditionObj := ormCondition.And("user_login", username).Or("user_email", email)
	var user models.UserModel
	cnt, err := o.QueryTable(user).SetCond(ormConditionObj).Count()
	if err != nil {
		logs.Error(err)
		return true
	} else if cnt > 0 {
		return true
	} else {
		return false
	}
}
