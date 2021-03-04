package main

import (
	"fmt"
	"github.com/beego/beego/v2/client/orm"
	"github.com/beego/beego/v2/core/logs"
	_ "github.com/go-sql-driver/mysql"
	_ "go-admin/routers"
	"go-admin/utils"
	"log"
	"os"
	"path/filepath"
	"time"

	beego "github.com/beego/beego/v2/server/web"
)

func init() {
	// 初始化日志
	initLogsDriver()
	// 初始化ORM
	initOrmDriver()
	// 获取当前运行目录
	utils.GetAppPath()
	//初始化session
	beego.BConfig.WebConfig.Session.SessionOn = true
	//跨域处理
	utils.CORS()
	// 文件上传初始化
	initFileUpdateDriver()
}

func main() {
	if beego.BConfig.RunMode == "dev" {
		beego.BConfig.WebConfig.DirectoryIndex = true
		beego.BConfig.WebConfig.StaticDir["/swagger"] = "swagger"
	}
	beego.Run()
}

// 初始化ORM
func initOrmDriver() {
	dbUser, _ := beego.AppConfig.String("db::dbUser")
	dbPass, _ := beego.AppConfig.String("db::dbPass")
	dbHost, _ := beego.AppConfig.String("db::dbHost")
	dbPort, _ := beego.AppConfig.String("db::dbPort")
	dbName, _ := beego.AppConfig.String("db::dbName")
	dbStr := dbUser + ":" + dbPass + "@tcp(" + dbHost + ":" + dbPort + ")/" + dbName + "?charset=utf8"
	// ORM日志
	orm.Debug = true
	// 输出到文件
	f, err := os.OpenFile("./logs/admin.orm.log", os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0644)
	if err != nil {
		logs.Error(err)
	}
	orm.DebugLog = orm.NewLog(f)
	// 注册driver为mysql
	registerDriverErr := orm.RegisterDriver("mysql", orm.DRMySQL)
	if registerDriverErr != nil {
		logs.Error("ORM registration failed." + registerDriverErr.Error())
		return
	}
	RegisterDataBaseErr := orm.RegisterDataBase("default", "mysql", dbStr)
	if RegisterDataBaseErr != nil {
		logs.Error("Register DataBase failed." + RegisterDataBaseErr.Error())
		return
	}
	// 设置为 UTC 时间
	orm.DefaultTimeLoc = time.UTC
}

// 初始化日志
func initLogsDriver() {
	if !utils.FileExists(utils.AppPath + "/logs") {
		// 创建日志目录
		err := os.MkdirAll(utils.AppPath + "/logs", os.ModePerm)
		if err != nil {
			fmt.Println(err)
			os.Exit(-1) //创建目录失败则终止程序
		}
	}
	logsErr := logs.SetLogger(logs.AdapterMultiFile, `{"filename":"./logs/admin.log","separate":["emergency", "alert", "critical", "error", "warning", "notice", "info", "debug"]}`)
	if logsErr != nil {
		//日志模块初始化失败
		log.Println("Beego log module initialization failed.")
		return
	}
	logs.Async() //日志异步
}

// 初始化文件上传
func initFileUpdateDriver() {
	beego.BConfig.MaxMemory = 1 << 27
	utils.UploadPath = filepath.Join(utils.AppPath, "upload") // 拼接路径

	// 检查目录是否存在
	if !utils.FileExists(utils.UploadPath) {
		err := os.MkdirAll(utils.UploadPath, os.ModePerm)
		if err != nil {
			logs.Error(err)
			os.Exit(-1) //创建目录失败则终止程序
		}
	}
}
