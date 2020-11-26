package lib

import (
	"crypto/md5"
	"encoding/hex"
	"github.com/astaxie/beego/logs"
	"github.com/astaxie/beego/session"
	"github.com/astaxie/beego/validation"
	"go-admin/models/admin"
	"os"
	"reflect"
	"unsafe"
)

//全局session
var (
	GlobalSessions *session.Manager
	CurrentUser    LoginUser
	AppPath        string //App运行目录
	UploadPath     string // 文件上传路径
)

type LoginUser struct {
	admin.UserModel
	IsRoot bool
}

// 密码加密
// @Param str String 待加密的字符串
// @return string
func Encryption(str string) string {
	var appKey = "Bd9JNSMx9VyBRX,lho7z0gVRgyBD5f!9"
	mdHash := md5.New()
	mdHash.Write([]byte(str))
	md5Byte := mdHash.Sum(nil)
	md5Str := hex.EncodeToString(md5Byte)
	rs := []rune(md5Str)          //把md5字符串转换成切片
	start := string(rs[0:6])      //开头截取6位
	end := string(rs[22:])        //结尾截取10位
	md5Str = start + appKey + end //加盐拼接
	mdHash.Write([]byte(md5Str))
	md5Byte = mdHash.Sum(nil)
	md5Str = hex.EncodeToString(md5Byte)
	return "###" + md5Str
}

func SessionInit() {
	sessionConfig := &session.ManagerConfig{
		CookieName:      "gosessionid",
		EnableSetCookie: true,
		Gclifetime:      3600,
		Maxlifetime:     3600,
		Secure:          false,
		CookieLifeTime:  3600,
		ProviderConfig:  "./tmp",
	}
	GlobalSessions, _ = session.NewManager("memory", sessionConfig)
	go GlobalSessions.GC()
}

//验证表单数据
func FormValidation(validData interface{}) (bool, string) {
	valid := validation.Validation{}
	b, err := valid.Valid(validData)
	if err != nil {
		// handle error
		logs.Error(err.Error())
		return false, err.Error()
	}

	//结果验证
	if !b {
		for _, err := range valid.Errors {
			msg := err.Field + " " + err.Message
			logs.Info(err.Key, err.Message)
			return false, msg
		}
	}
	return true, ""
}

func BytesToString(b []byte) string {
	bh := (*reflect.SliceHeader)(unsafe.Pointer(&b))
	sh := reflect.StringHeader{bh.Data, bh.Len}
	return *(*string)(unsafe.Pointer(&sh))
}

func StringToBytes(s string) []byte {
	sh := (*reflect.StringHeader)(unsafe.Pointer(&s))
	bh := reflect.SliceHeader{sh.Data, sh.Len, 0}
	return *(*[]byte)(unsafe.Pointer(&bh))
}

// 判断文件或者目录是否存在
func FileExists(filename string) bool {
	_, err := os.Stat(filename)
	return err == nil || os.IsExist(err)
}
