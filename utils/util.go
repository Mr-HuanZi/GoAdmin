package utils

import (
	"go-admin/models/admin"
	"log"
	"os"
	"path/filepath"
)

// 全局变量

var (
	CurrentUser    LoginUser
	AppPath        string //App运行目录
	UploadPath     string // 文件上传路径
)

type LoginUser struct {
	admin.UserModel
	IsRoot bool
}

// 获取当前运行目录的绝对路径
func GetAppPath() {
	dir, err := filepath.Abs(filepath.Dir(os.Args[0]))
	if err != nil {
		log.Fatal(err)
	}
	AppPath = dir
}

