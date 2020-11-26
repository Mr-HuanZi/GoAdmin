package admin

import (
	"crypto/md5"
	"encoding/hex"
	"github.com/astaxie/beego/logs"
	"github.com/astaxie/beego/orm"
	"go-admin/lib"
	"go-admin/models/admin"
	"io"
	"os"
	"path"
	"path/filepath"
	"strconv"
	"time"
)

type FileUploadController struct {
	BaseController
}

// 单文件上传
func (c *FileUploadController) UploadFile() {
	f, h, err := c.GetFile("file")

	if err != nil {
		logs.Error("GetFile err ", err)
		c.Response(500, "", nil)
	}
	logs.Debug(h.Header, h.Filename, h.Size)
	// 关闭文件
	defer func() {
		err := f.Close()
		// 关闭文件错误
		if err != nil {
			logs.Error("File close err ", err)
		}
	}()

	// 文件信息
	var asset admin.AssetModel
	asset.Status = 1
	asset.FileInfo = h.Filename
	asset.Size = h.Size

	// 创建文件保存目录
	savePath := filepath.Join(lib.UploadPath, time.Now().Format("2006-01-02")) // 绝对路径
	if !lib.FileExists(savePath) {
		err := os.MkdirAll(savePath, os.ModePerm)
		if err != nil {
			logs.Error("savePath make err ", err)
			c.Response(500, "", nil)
			return
		}
	}

	// 对文件进行MD5
	mdHash := md5.New()
	_, _ = io.Copy(mdHash, f) // f所在的缓存区内容将被清空
	MD5Str := hex.EncodeToString(mdHash.Sum(nil))
	asset.Md5 = MD5Str

	// 获得文件名
	// 组合字符串组成MD5字串，避免文件重名
	mdHash.Write([]byte(MD5Str + h.Filename + strconv.FormatInt(h.Size, 10) + strconv.FormatInt(time.Now().UnixNano()/1e6, 10)))
	MD5Str = hex.EncodeToString(mdHash.Sum(nil))
	filename := MD5Str + path.Ext(h.Filename) // 获取文件后缀
	asset.FileName = filename
	// 文件保存的绝对路径
	filePath := filepath.Join(savePath, filename)
	logs.Debug(filePath)
	logs.Debug(lib.AppPath)
	relPath, relErr := filepath.Rel(lib.AppPath, filePath)
	if relErr != nil {
		logs.Error(relErr)
		c.Response(500, "", nil)
		return
	}
	logs.Debug(filepath.Rel(lib.AppPath, relPath))
	// 相对路径
	asset.Path = relPath

	// 保存文件
	// 由于md5加密的时候，已经将f从缓冲区读取完毕，因此需要重新载入缓冲
	// 所以直接调用SaveToFile方法即可
	_ = c.SaveToFile("file", filePath)

	asset.CreateTime = time.Now().UnixNano() / 1e6
	asset.AddStaff = lib.CurrentUser.Id

	o := orm.NewOrm()
	id, insertErr := o.Insert(&asset)
	if insertErr != nil {
		logs.Error(insertErr)
		c.Response(500, "", nil)
	} else {
		asset.Id = id
		c.Response(200, "", asset)
	}
}
