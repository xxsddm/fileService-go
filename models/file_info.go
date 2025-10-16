package models

import (
	"time"
)

// FileInfo 文件信息实体
type FileInfo struct {
	ID         uint64    `gorm:"primaryKey;column:id" json:"id"`
	FileName   string    `gorm:"column:fileName" json:"fileName"`
	FilePath   string    `gorm:"column:filePath" json:"filePath"`
	FileSize   int       `gorm:"column:fileSize" json:"fileSize"`
	Status     int       `gorm:"column:status" json:"status"`
	UploadDate time.Time `gorm:"column:uploadDate" json:"uploadDate"`
}

// TableName 设置表名
func (FileInfo) TableName() string {
	return "fileInfo"
}

// FileInfoDTO 文件信息DTO
type FileInfoDTO struct {
	ID         uint64    `json:"id"`
	FileName   string    `json:"fileName"`
	FilePath   string    `json:"filePath"`
	FileSize   int       `json:"fileSize"`
	Status     int       `json:"status"`
	UploadDate time.Time `json:"uploadDate"`
}
