package lib

import (
	"crypto/md5"
	"encoding/hex"
	"errors"
	"github.com/beego/beego/v2/core/logs"
	"github.com/beego/beego/v2/core/validation"
	"go-admin/models/admin"
	"io"
	"mime/multipart"
	"os"
	"reflect"
	"unsafe"
)

//全局session
var (
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
	md5Str := Md5(str)
	rs := []rune(md5Str)          //把md5字符串转换成切片
	start := string(rs[0:6])      //开头截取6位
	end := string(rs[22:])        //结尾截取10位
	md5Str = start + appKey + end //加盐拼接
	md5Str = Md5(md5Str)
	return "###" + md5Str
}

// MD5加密封装
func Md5(str string) string {
	mdHash := md5.New()
	mdHash.Write([]byte(str))
	md5Byte := mdHash.Sum(nil)
	return hex.EncodeToString(md5Byte)
}

// MD5文件加密封装
// 注意，调用此方法可能会导致multipart.File所在的缓冲区被清空
func Md5File(f multipart.File) string {
	mdHash := md5.New()
	_, _ = io.Copy(mdHash, f) // f所在的缓存区内容将被清空
	return hex.EncodeToString(mdHash.Sum(nil))
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

// 把结构体right的值复制给left
// left interface{} 被修改的结构体
// right interface{} 有数据的结构体
func StructCopy(left interface{}, right interface{}) error {
	lVal := reflect.ValueOf(left)
	lValElem := lVal.Elem()
	lValType := reflect.TypeOf(left)

	rVal := reflect.ValueOf(right)
	rValElem := rVal.Elem()
	rValType := reflect.TypeOf(right)

	if lValType.Kind() != reflect.Ptr || lValType.Elem().Kind() == reflect.Ptr || rValType.Kind() != reflect.Ptr || rValType.Elem().Kind() == reflect.Ptr {
		return errors.New("type of parameters must be Ptr of value")
	}

	if lVal.IsNil() || rVal.IsNil() {
		return errors.New("value of parameters should not be nil")
	}

	rTypeOfT := rValElem.Type()

	for i := 0; i < rValElem.NumField(); i++ {
		r := rTypeOfT.Field(i)
		if r.Anonymous {
			continue // 跳过嵌套结构体
		}
		lField := lValElem.FieldByName(r.Name)
		rField := rValElem.FieldByName(r.Name)
		if !lField.IsValid() {
			continue
		}

		// 类型不相等不能赋值
		if lField.Type() != rField.Type() {
			continue
		}
		// 在要修改的结构体中查询有数据结构体中相同属性的字段，有则修改其值
		lField.Set(reflect.ValueOf(rValElem.Field(i).Interface()))
	}
	return nil
}

// 位移标志
// 用于单字段表示多状态时的状态改变
// isWhat bool 启用/禁用 状态
// bit int 状态预设值
// flag int 当前状态
// return int 改变后的flag
func ShiftFlag(isWhat bool, bit int, flag int) int {
	logs.Debug("*********************ShiftFlag*********************")
	logs.Debug("%b", bit)
	if isWhat {
		flag |= bit
		logs.Debug("%b", flag)
	} else {
		// 以下方法不保证准确性，还在测试阶段
		flag |= bit // 先按位或，得到相反的值
		logs.Debug("%b", flag)
		flag ^= bit // 再按位异或，得到最终结果
		logs.Debug("%b", flag)
	}
	logs.Debug("%b", flag)
	logs.Debug("---------------------ShiftFlag---------------------")
	return flag
}
