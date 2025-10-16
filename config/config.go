package config

import (
	"github.com/spf13/viper"
)

// Config 全局配置结构体
type Config struct {
	Server   ServerConfig   `mapstructure:"server"`
	Database DatabaseConfig `mapstructure:"database"`
	File     FileConfig     `mapstructure:"file"`
}

// ServerConfig 服务器配置
type ServerConfig struct {
	Port int `mapstructure:"port"`
}

// DatabaseConfig 数据库配置
type DatabaseConfig struct {
	Host     string `mapstructure:"host"`
	Port     int    `mapstructure:"port"`
	User     string `mapstructure:"user"`
	Password string `mapstructure:"password"`
	Name     string `mapstructure:"name"`
	Charset  string `mapstructure:"charset"`
}

// FileConfig 文件配置
type FileConfig struct {
	Upload UploadConfig `mapstructure:"upload"`
}

// UploadConfig 上传配置
type UploadConfig struct {
	Path         string `mapstructure:"path"`
	MaxSize      int64  `mapstructure:"max_size"`
	AllowedTypes string `mapstructure:"allowed_types"`
}

var AppConfig *Config

// LoadConfig 加载配置
func LoadConfig() error {
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	viper.AddConfigPath(".")

	// 设置默认值
	viper.SetDefault("server.port", 8080)
	viper.SetDefault("database.host", "localhost")
	viper.SetDefault("database.port", 3306)
	viper.SetDefault("database.user", "root")
	viper.SetDefault("database.password", "123456")
	viper.SetDefault("database.name", "file_service")
	viper.SetDefault("database.charset", "utf8mb4")
	viper.SetDefault("file.upload.path", "./uploads")
	viper.SetDefault("file.upload.max_size", 104857600) // 100MB
	viper.SetDefault("file.upload.allowed_types", "jpg,jpeg,png,gif,pdf,doc,docx,xls,xlsx,txt,zip,rar")

	if err := viper.ReadInConfig(); err != nil {
		return err
	}

	AppConfig = &Config{}
	if err := viper.Unmarshal(AppConfig); err != nil {
		return err
	}

	return nil
}
