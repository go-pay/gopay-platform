package router

import (
	"github.com/gin-gonic/gin"
	"github.com/go-pay/web"
	"github.com/go-pay/xlog"
)

// 文件上传
func fileUpload(c *gin.Context) {
	contentType := c.GetHeader("Content-Type")
	xlog.Infof("收到文件上传请求, Content-Type: %s", contentType)

	// 尝试 multipart/form-data 上传
	file, err := c.FormFile("file")
	if err == nil {
		// 表单上传成功
		rsp, err := svc.FileUpload(c, file)
		if err != nil {
			web.JSON(c, nil, err)
			return
		}
		web.JSON(c, rsp, nil)
		return
	}

	// 如果表单获取失败，尝试二进制流上传
	xlog.Infof("表单获取失败，尝试二进制流上传: %v", err)
	rsp, err := svc.FileUploadBinary(c)
	if err != nil {
		xlog.Errorf("二进制流上传失败: %v", err)
		web.JSON(c, nil, err)
		return
	}
	web.JSON(c, rsp, nil)
}
