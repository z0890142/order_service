package config

import (
	"sync"
)

var globalConfig *Config
var configOnce sync.Once

// ResetConfig set config to Nil, used for tests
func ResetConfig() {
	globalConfig = nil
}

// GetConfig 獲取該服務相關配置
func GetConfig() *Config {
	configOnce.Do(func() {
		globalConfig = &Config{}
	})
	return globalConfig
}

// Config 該服務相關配置
type Config struct {
	Env     string  `mapstructure:"ENV"`
	Service Service `mapstructure:"SERVICE"`

	LogLevel     string   `mapstructure:"LOG_LEVEL"`
	LogFile      []string `mapstructure:"LOG_FILE"`
	ErrorLogFile []string `mapstructure:"ERROR_LOG_FILE"`

	Databases         DatabaseOption `mapstructure:"DATABASES"`
	Mongo             Mongo          `mapstructure:"MONGO"`
	MigrationFilePath string         `mapstructure:"MIGRATION_FILE_PATH"`
}

type Mongo struct {
	Url        string `mapstructure:"URL"`
	Collection string `mapstructure:"COLLECTION"`
	Database   string `mapstructure:"DATABASE"`
}

type DatabaseOption struct {
	Driver   string `mapstructure:"DRIVER"`
	Host     string `mapstructure:"HOST"`
	Port     uint16 `mapstructure:"PORT"`
	Username string `mapstructure:"USERNAME"`
	Password string `mapstructure:"PASSWORD"`
	DBName   string `mapstructure:"DBNAME"`
	PoolSize int    `mapstructure:"POOL_SIZE"`
}

type Service struct {
	Name string `mapstructure:"NAME"`
	Host string `mapstructure:"HOST"`
	Port string `mapstructure:"PORT"`
}
