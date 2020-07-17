package main

import (
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/logs"
	"github.com/astaxie/beego/orm"
	"github.com/astaxie/beego/plugins/cors"
	_ "github.com/go-sql-driver/mysql"
	"go-admin/lib"
	_ "go-admin/routers"
	"log"
	"time"
)

func init() {
	// 初始化ORM
	initOrmDriver()
	//设置日志
	initLogsDriver()
	//初始化session
	beego.BConfig.WebConfig.Session.SessionOn = true
	lib.SessionInit()
	//跨域处理
	beego.InsertFilter("*", beego.BeforeRouter, cors.Allow(&cors.Options{
		//允许访问的源
		AllowOrigins: []string{"http://localhost"},
		//允许访问所有源
		//AllowAllOrigins: true,
		//可选参数"GET", "POST", "PUT", "DELETE", "OPTIONS" (*为所有)
		//其中Options跨域复杂请求预检
		AllowMethods: []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		//指的是允许的Header的种类
		AllowHeaders: []string{"Origin", "Authorization", "Access-Control-Allow-Origin", "Access-Control-Allow-Headers", "Content-Type"},
		//公开的HTTP标头列表
		ExposeHeaders: []string{"Content-Length", "Access-Control-Allow-Origin", "Access-Control-Allow-Headers", "Content-Type"},
		//如果设置，则允许共享身份验证凭据，例如cookie
		AllowCredentials: true,
	}))
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
	dbUser := beego.AppConfig.String("db::dbUser")
	dbPass := beego.AppConfig.String("db::dbPass")
	dbHost := beego.AppConfig.String("db::dbHost")
	dbPort := beego.AppConfig.String("db::dbPort")
	dbName := beego.AppConfig.String("db::dbName")
	dbStr := dbUser + ":" + dbPass + "@tcp(" + dbHost + ":" + dbPort + ")/" + dbName + "?charset=utf8"
	orm.Debug = true
	ormErr := orm.RegisterDriver("mysql", orm.DRMySQL)
	if ormErr != nil {
		log.Println("ORM registration failed.")
		return
	}
	ormErr = orm.RegisterDataBase("default", "mysql", dbStr)
	if ormErr != nil {
		log.Println("Register DataBase failed.")
		return
	}
	// 设置为 UTC 时间
	orm.DefaultTimeLoc = time.UTC
}

// 初始化日志
func initLogsDriver() {
	logsErr := logs.SetLogger(logs.AdapterMultiFile, `{"filename":"./logs/admin.log","separate":["emergency", "alert", "critical", "error", "warning", "notice", "info", "debug"]}`)
	if logsErr != nil {
		//日志模块初始化失败
		log.Println("Beego log module initialization failed.")
		return
	}
	logs.Async() //日志异步
}
