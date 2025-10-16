package controllers

import (
	"fmt"
	"net/http"
	"net/url"
	"strconv"
	"file-service/service"

	"github.com/gin-gonic/gin"
)

// FileController 文件控制器
type FileController struct {
	fileService *service.FileService
}

// NewFileController 创建文件控制器实例
func NewFileController() *FileController {
	return &FileController{
		fileService: &service.FileService{},
	}
}

// UploadFiles 批量上传文件
func (fc *FileController) UploadFiles(c *gin.Context) {
	// 解析表单
	if err := c.Request.ParseMultipartForm(32 << 20); err != nil { // 32MB
		c.JSON(http.StatusBadRequest, gin.H{"error": "failed to parse form"})
		return
	}

	// 获取文件
	form, err := c.MultipartForm()
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "failed to get multipart form"})
		return
	}

	files := form.File["files"]
	if len(files) == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "no files provided"})
		return
	}

	// 调用服务上传文件
	fileInfos, err := fc.fileService.UploadFiles(files)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": fmt.Sprintf("upload failed: %v", err),
			"files":   nil,
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "upload success",
		"files":   fileInfos,
	})
}

// DownloadFile 下载文件
func (fc *FileController) DownloadFile(c *gin.Context) {
	fileName, err := url.QueryUnescape(c.Param("fileName"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid file name"})
		return
	}

	fileData, err := fc.fileService.DownloadFile(fileName)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "file not found"})
		return
	}

	// 设置响应头
	c.Header("Content-Type", "application/octet-stream")
	c.Header("Content-Disposition", fmt.Sprintf("attachment; filename=\"%s\"", fileName))
	c.Data(http.StatusOK, "application/octet-stream", fileData)
}

// GetFileList 获取文件列表
func (fc *FileController) GetFileList(c *gin.Context) {
	// 解析查询参数
	page, err := strconv.Atoi(c.DefaultQuery("page", "1"))
	if err != nil || page < 1 {
		page = 1
	}

	pageSize, err := strconv.Atoi(c.DefaultQuery("page_size", "10"))
	if err != nil || pageSize < 1 {
		pageSize = 10
	}
	if pageSize > 100 {
		pageSize = 100 // 限制最大页面大小
	}

	fileNameFilter := c.Query("filename_filter")

	// 调用服务获取文件列表
	result, err := fc.fileService.GetFileList(page, pageSize, fileNameFilter)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": fmt.Sprintf("query file list failed: %v", err),
			"files":   nil,
		})
		return
	}

	c.JSON(http.StatusOK, result)
}