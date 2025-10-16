package service

import (
	"file-service/config"
	"file-service/database"
	"file-service/dto"
	"file-service/models"
	"file-service/utils"
	"fmt"
	"io"
	"mime/multipart"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/jinzhu/gorm"
)

// FileService 文件服务接口
type FileService struct{}

// UploadFiles 批量上传文件
func (fs *FileService) UploadFiles(files []*multipart.FileHeader) ([]models.FileInfoDTO, error) {
	var results []models.FileInfoDTO
	var fileInfos []models.FileInfo

	uploadPath := config.AppConfig.File.Upload.Path

	// 确保上传目录存在
	if err := os.MkdirAll(uploadPath, os.ModePerm); err != nil {
		return nil, fmt.Errorf("failed to create upload directory: %v", err)
	}

	for _, fileHeader := range files {
		// 验证文件
		if err := fs.validateFile(fileHeader); err != nil {
			return nil, fmt.Errorf("file validation failed: %v", err)
		}

		// 保存文件
		fileName, err := fs.saveFile(fileHeader)
		if err != nil {
			return nil, fmt.Errorf("failed to save file: %v", err)
		}

		// 创建文件信息记录
		fileInfo := models.FileInfo{
			ID:         utils.GenerateID(),
			FileName:   fileName,
			FilePath:   filepath.Join(uploadPath, fileName),
			FileSize:   int(fileHeader.Size),
			Status:     0,
			UploadDate: time.Now(),
		}

		// 保存到数据库
		if err := database.DB.Create(&fileInfo).Error; err != nil {
			return nil, fmt.Errorf("failed to save file info to database: %v", err)
		}

		fileInfos = append(fileInfos, fileInfo)
		results = append(results, fs.convertToDTO(fileInfo))
	}

	return results, nil
}

// DownloadFile 下载文件
func (fs *FileService) DownloadFile(fileName string) ([]byte, error) {
	var fileInfo models.FileInfo

	// 从数据库查询文件信息
	if err := database.DB.Where("fileName = ? AND status = ?", fileName, 0).First(&fileInfo).Error; err != nil {
		if gorm.IsRecordNotFoundError(err) {
			return nil, fmt.Errorf("file not found: %s", fileName)
		}
		return nil, fmt.Errorf("database query error: %v", err)
	}

	// 检查文件是否存在
	if _, err := os.Stat(fileInfo.FilePath); os.IsNotExist(err) {
		// 文件不存在，标记为无效
		database.DB.Model(&fileInfo).Update("status", 1)
		return nil, fmt.Errorf("file not found on disk: %s", fileName)
	}

	// 读取文件内容
	fileData, err := os.ReadFile(fileInfo.FilePath)
	if err != nil {
		return nil, fmt.Errorf("failed to read file: %v", err)
	}

	return fileData, nil
}

// GetFileList 分页查询文件列表
func (fs *FileService) GetFileList(page, pageSize int, fileNameFilter string) (*dto.PageResult[models.FileInfoDTO], error) {
	offset := (page - 1) * pageSize

	var fileInfos []models.FileInfo
	var total int

	// 构建查询
	query := database.DB.Where("status = ?", 0)

	// 添加文件名过滤条件
	if fileNameFilter != "" {
		query = query.Where("fileName LIKE ?", "%"+fileNameFilter+"%")
	}

	// 查询总数
	query.Model(&models.FileInfo{}).Count(&total)

	// 分页查询
	query.Order("uploadDate DESC").Offset(offset).Limit(pageSize).Find(&fileInfos)

	// 转换为DTO
	var fileInfoDTOs []models.FileInfoDTO
	for _, fileInfo := range fileInfos {
		fileInfoDTOs = append(fileInfoDTOs, fs.convertToDTO(fileInfo))
	}

	// 创建分页结果
	pageResult := dto.NewPageResult(fileInfoDTOs, total, page, pageSize)

	return pageResult, nil
}

// DeleteFilesByIds 根据ID批量删除文件
func (fs *FileService) DeleteFilesByIds(ids []uint64) (int, error) {
	if len(ids) == 0 {
		return 0, nil
	}

	// 分批处理，每批500个
	batchSize := 500
	count := 0

	for i := 0; i < len(ids); i += batchSize {
		end := i + batchSize
		if end > len(ids) {
			end = len(ids)
		}

		batchIds := ids[i:end]

		// 批量查询有效文件
		var fileInfos []models.FileInfo
		if err := database.DB.Where("id IN (?) AND status = ?", batchIds, 0).Find(&fileInfos).Error; err != nil {
			return count, fmt.Errorf("failed to query files: %v", err)
		}

		var validIds []uint64

		// 删除物理文件
		for _, fileInfo := range fileInfos {
			if err := os.Remove(fileInfo.FilePath); err != nil && !os.IsNotExist(err) {
				// 记录错误但继续处理其他文件
				continue
			}
			validIds = append(validIds, fileInfo.ID)
		}

		// 批量更新数据库状态
		if len(validIds) > 0 {
			if err := database.DB.Table("fileInfo").Where("id IN (?)", validIds).Update("status", 1).Error; err != nil {
				return count, fmt.Errorf("failed to update file status: %v", err)
			}
			count += len(validIds)
		}
	}

	return count, nil
}

// CleanExpiredFiles 清理过期文件（7天前的文件）
func (fs *FileService) CleanExpiredFiles() {
	// 计算7天前的时间
	expiryDate := time.Now().AddDate(0, 0, -7)

	var expiredFiles []models.FileInfo

	// 查询7天前的文件
	if err := database.DB.Where("status = ? AND uploadDate < ?", 0, expiryDate).Find(&expiredFiles).Error; err != nil {
		return
	}

	for _, fileInfo := range expiredFiles {
		// 删除物理文件
		if err := os.Remove(fileInfo.FilePath); err != nil && !os.IsNotExist(err) {
			continue
		}

		// 更新数据库状态
		database.DB.Model(&fileInfo).Update("status", 1)
	}
}

// validateFile 验证文件
func (fs *FileService) validateFile(fileHeader *multipart.FileHeader) error {
	// 检查文件是否为空
	if fileHeader.Size == 0 {
		return fmt.Errorf("file is empty")
	}

	// 检查文件大小
	maxSize := config.AppConfig.File.Upload.MaxSize
	if fileHeader.Size > maxSize {
		return fmt.Errorf("file size exceeds the limit of %d bytes", maxSize)
	}

	// 检查文件名
	if fileHeader.Filename == "" {
		return fmt.Errorf("file name is empty")
	}

	// 检查文件类型
	fileExtension := strings.ToLower(filepath.Ext(fileHeader.Filename))
	if fileExtension != "" {
		fileExtension = fileExtension[1:] // 移除点号
	}

	allowedTypes := strings.Split(config.AppConfig.File.Upload.AllowedTypes, ",")
	allowed := false
	for _, allowedType := range allowedTypes {
		if strings.TrimSpace(allowedType) == fileExtension {
			allowed = true
			break
		}
	}

	if !allowed {
		return fmt.Errorf("file type %s is not supported", fileExtension)
	}

	return nil
}

// saveFile 保存文件
func (fs *FileService) saveFile(fileHeader *multipart.FileHeader) (string, error) {
	// 打开上传的文件
	file, err := fileHeader.Open()
	if err != nil {
		return "", fmt.Errorf("failed to open uploaded file: %v", err)
	}
	defer file.Close()

	// 生成唯一文件名
	originalFileName := fileHeader.Filename
	fileExtension := filepath.Ext(originalFileName)
	baseName := strings.TrimSuffix(originalFileName, fileExtension)
	uuidStr := uuid.New().String()
	fileName := fmt.Sprintf("%s_%s%s", baseName, uuidStr, fileExtension)

	// 构建文件路径
	uploadPath := config.AppConfig.File.Upload.Path
	filePath := filepath.Join(uploadPath, fileName)

	// 创建目标文件
	dst, err := os.Create(filePath)
	if err != nil {
		return "", fmt.Errorf("failed to create destination file: %v", err)
	}
	defer dst.Close()

	// 复制文件内容
	if _, err := io.Copy(dst, file); err != nil {
		return "", fmt.Errorf("failed to copy file content: %v", err)
	}

	return fileName, nil
}

// convertToDTO 转换为DTO
func (fs *FileService) convertToDTO(fileInfo models.FileInfo) models.FileInfoDTO {
	return models.FileInfoDTO{
		ID:         fileInfo.ID,
		FileName:   fileInfo.FileName,
		FilePath:   fileInfo.FilePath,
		FileSize:   fileInfo.FileSize,
		Status:     fileInfo.Status,
		UploadDate: fileInfo.UploadDate,
	}
}
