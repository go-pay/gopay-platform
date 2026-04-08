package router

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	"gopay/errcode"

	"github.com/gin-gonic/gin"
	"github.com/go-pay/web"
	"github.com/go-pay/xlog"
	"github.com/google/uuid"
)

const (
	maxUploadSize = 5 << 20 // 5MB
	uploadDir     = "./uploads"
)

func uploadImage(c *gin.Context) {
	file, err := c.FormFile("file")
	if err != nil {
		xlog.Errorf("uploadImage FormFile, err:%v", err)
		web.JSON(c, nil, errcode.RequestErr)
		return
	}
	// 校验文件大小
	if file.Size > maxUploadSize {
		web.JSON(c, nil, errcode.RequestErr)
		return
	}
	// 校验文件类型
	contentType := file.Header.Get("Content-Type")
	if !strings.HasPrefix(contentType, "image/") {
		web.JSON(c, nil, errcode.RequestErr)
		return
	}
	// 生成存储路径: uploads/YYYY/MM/uuid.ext
	ext := filepath.Ext(file.Filename)
	now := time.Now()
	dir := filepath.Join(uploadDir, now.Format("2006"), now.Format("01"))
	if err = os.MkdirAll(dir, 0755); err != nil {
		xlog.Errorf("uploadImage MkdirAll, err:%v", err)
		web.JSON(c, nil, errcode.ServerErr)
		return
	}
	filename := fmt.Sprintf("%s%s", uuid.New().String(), ext)
	savePath := filepath.Join(dir, filename)
	if err = c.SaveUploadedFile(file, savePath); err != nil {
		xlog.Errorf("uploadImage SaveUploadedFile, err:%v", err)
		web.JSON(c, nil, errcode.ServerErr)
		return
	}
	// 返回相对 URL
	url := fmt.Sprintf("/uploads/%s/%s/%s", now.Format("2006"), now.Format("01"), filename)
	web.JSON(c, gin.H{"url": url}, nil)
}
