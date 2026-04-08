package service

import (
	"fmt"
	"io"
	"mime/multipart"
	"os"
	"path/filepath"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/go-pay/xlog"
)

// FileUploadReq 文件上传请求
type FileUploadReq struct {
	File *multipart.FileHeader `form:"file" binding:"required"`
}

// FileUploadRsp 文件上传响应
type FileUploadRsp struct {
	FileName string `json:"file_name"`
	FilePath string `json:"file_path"`
	FileSize int64  `json:"file_size"`
}

// FileUploadBinary 二进制流文件上传
func (s *Service) FileUploadBinary(c *gin.Context) (*FileUploadRsp, error) {
	// 从 query 参数或 header 获取文件名
	fileName := c.Query("filename")
	if fileName == "" {
		fileName = c.GetHeader("X-Filename")
	}
	if fileName == "" {
		// 尝试从 Content-Disposition 获取
		contentDisposition := c.GetHeader("Content-Disposition")
		if contentDisposition != "" {
			// 简单解析 filename
			if idx := filepath.Base(contentDisposition); idx != "" {
				fileName = idx
			}
		}
	}
	if fileName == "" {
		// 根据 Content-Type 推断扩展名
		contentType := c.GetHeader("Content-Type")
		ext := getExtByContentType(contentType)
		fileName = fmt.Sprintf("%s%s", time.Now().Format("20060102150405"), ext)
	}

	// 生成唯一文件名
	ext := filepath.Ext(fileName)
	timestamp := time.Now().Format("20060102150405")
	newFileName := fmt.Sprintf("%s_%s%s", timestamp, filepath.Base(fileName[:len(fileName)-len(ext)]), ext)

	// 目标路径
	uploadDir := "file"
	if err := os.MkdirAll(uploadDir, 0755); err != nil {
		xlog.Errorf("创建目录失败: %v", err)
		return nil, err
	}

	filePath := filepath.Join(uploadDir, newFileName)

	// 保存文件
	dst, err := os.Create(filePath)
	if err != nil {
		xlog.Errorf("创建目标文件失败: %v", err)
		return nil, err
	}
	defer dst.Close()

	written, err := io.Copy(dst, c.Request.Body)
	if err != nil {
		xlog.Errorf("保存文件失败: %v", err)
		return nil, err
	}

	return &FileUploadRsp{
		FileName: newFileName,
		FilePath: filePath,
		FileSize: written,
	}, nil
}

// FileUpload 文件上传
func (s *Service) FileUpload(c *gin.Context, file *multipart.FileHeader) (*FileUploadRsp, error) {
	// 生成文件名（使用时间戳避免重复）
	ext := filepath.Ext(file.Filename)
	timestamp := time.Now().Format("20060102150405")
	newFileName := fmt.Sprintf("%s_%s%s", timestamp, filepath.Base(file.Filename[:len(file.Filename)-len(ext)]), ext)

	// 目标路径
	uploadDir := "file"
	if err := os.MkdirAll(uploadDir, 0755); err != nil {
		xlog.Errorf("创建目录失败: %v", err)
		return nil, err
	}

	filePath := filepath.Join(uploadDir, newFileName)

	// 保存文件
	src, err := file.Open()
	if err != nil {
		xlog.Errorf("打开上传文件失败: %v", err)
		return nil, err
	}
	defer src.Close()

	dst, err := os.Create(filePath)
	if err != nil {
		xlog.Errorf("创建目标文件失败: %v", err)
		return nil, err
	}
	defer dst.Close()

	if _, err = io.Copy(dst, src); err != nil {
		xlog.Errorf("保存文件失败: %v", err)
		return nil, err
	}

	return &FileUploadRsp{
		FileName: newFileName,
		FilePath: filePath,
		FileSize: file.Size,
	}, nil
}

// getExtByContentType 根据 Content-Type 获取文件扩展名
func getExtByContentType(contentType string) string {
	extMap := map[string]string{
		"application/zip":              ".zip",
		"application/x-zip-compressed": ".zip",
		"application/pdf":              ".pdf",
		"image/jpeg":                   ".jpg",
		"image/png":                    ".png",
		"image/gif":                    ".gif",
		"application/json":             ".json",
		"text/plain":                   ".txt",
		"application/octet-stream":     ".bin",
		"application/vnd.ms-excel":     ".xls",
		"application/vnd.openxmlformats-officedocument.spreadsheetml.sheet": ".xlsx",
		"application/msword": ".doc",
		"application/vnd.openxmlformats-officedocument.wordprocessingml.document": ".docx",
	}

	if ext, ok := extMap[contentType]; ok {
		return ext
	}
	return ".bin"
}
