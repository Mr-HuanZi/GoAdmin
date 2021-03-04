# GoAdmin
基于beego的后台管理系统

## 安装说明
> 需要开启 `Go Mod`

#### 1、下载项目
```bash
git clone git@gitee.com:Mr_Huan/go-admin.git
#或
git clone https://gitee.com/Mr_Huan/go-admin.git
```
或着使用go get
```bash
go get github.com/Mr-HuanZi/GoAdmin
```

#### 2、安装bee工具
```bash
go get -u github.com/beego/bee/v2
```

#### 3、安装依赖包
```bash
go mod vendor
```

#### 4、启动项目
```bash
bee run
```