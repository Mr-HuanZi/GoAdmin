package admin

import (
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/logs"
	"github.com/astaxie/beego/orm"
	"time"
)

type UserModel struct {
	Id                int64  `orm:"pk"`
	UserType          int8   `valid:"Range(0,2)"` //用户类别
	Sex               int8   //性别
	Birthday          int    //生日
	LastLoginTime     int64  //最后登录时间
	LastLoginIp       string //最后登录IP
	Score             int    //用户积分
	Gold              int    //金币
	CreateTime        int64  //注册时间
	UpdateTime        int64
	UserStatus        int8   //用户状态;0:禁用,1:正常,2:未验证
	UserLogin         string `valid:"Required;AlphaNumeric"` //用户名
	UserPass          string `valid:"Required;AlphaDash"`    //登录密码
	UserNickname      string `valid:"Required"`              //用户昵称
	UserEmail         string //用户登录邮箱
	UserUrl           string //用户个人网址
	Avatar            int    //用户头像
	Signature         string //个性签名
	UserActivationKey string //激活码
	Mobile            string //手机号
	Position          int    //层级
	LockTime          int    //登陆错误锁定开始时间
	ErrorSum          int8   //登陆错误次数
	First             int8   //是否首次登录系统
	LastEditPass      int    //最后一次修改密码的时间
	Openid            string //微信openid
}

func init() {
	//设置表前缀并且注册模型
	orm.RegisterModelWithPrefix(beego.AppConfig.String("db::dbPrefix"), new(UserModel))
}

//自定义表名
func (User *UserModel) TableName() string {
	return "user"
}

func Login(username string, password string) (int, int64) {
	o := orm.NewOrm()
	var user UserModel
	//cut,err := o.QueryTable(user).Filter("user_login", username).Filter("user_pass", password).Count()
	//查询完整的记录
	err := o.QueryTable(user).Filter("user_login", username).Filter("user_pass", password).One(&user)
	if err == orm.ErrMultiRows {
		// 查询到多个记录
		return 106, 0
	} else if user.Id != 0 {
		return 100, user.Id
	} else {
		return 102, 0
	}
}

//更新用户登录信息
func UpdateUserLoginInfo(uid int64, ip string) {
	o := orm.NewOrm()
	user := UserModel{Id: uid}
	if o.Read(&user) == nil {
		user.LastLoginTime = time.Now().Unix() //获取当前登录时间
		user.LastLoginIp = ip                  //记录登录IP
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

// 检查用户是否重复
func CheckUserRepeat(username string, email string) bool {
	o := orm.NewOrm()
	//自定义条件
	ormCondition := orm.NewCondition()
	ormConditionObj := ormCondition.And("user_login", username).Or("user_email", email)
	var user UserModel
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

//创建新用户
func CreateUser(userData *UserModel) int64 {
	o := orm.NewOrm()
	id, err := o.Insert(userData)
	if err == nil {
		logs.Info(id)
		return id
	}
	return 0
}

// 获取一个用户信息
func GetUser(uid int64) (UserModel, error) {
	o := orm.NewOrm()
	user := UserModel{Id: uid}
	err := o.Read(&user)
	if err == nil {
		return user, nil
	}
	return UserModel{}, err
}

// 改变用户的状态值
func ChangeUserStatus(uid []int64, status int8) (int64, error) {
	o := orm.NewOrm()
	return o.QueryTable(new(UserModel)).Filter("id__in", uid).Update(orm.Params{
		"user_status": status,
	})
}

// 检查用户名是否重复
func CheckUserDuplication(username string) bool {
	o := orm.NewOrm()
	// 检查用户名是否存在
	count, countErr := o.QueryTable(new(UserModel)).Filter("user_login", username).Count()
	if countErr != nil {
		return true
	}
	if count > 0 {
		logs.Notice("用户名[", username, "]", "已存在[", count, "]个")
		return true
	}
	return false
}
