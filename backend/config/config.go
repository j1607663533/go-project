package config

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	Host string
	Port string
	Mode string
	DB   DatabaseConfig
	AI   AIConfig
}

type DatabaseConfig struct {
	Host     string
	Port     string
	User     string
	Password string
	DBName   string
}

type AIConfig struct {
	APIKey  string
	BaseURL string
	Model   string
}

var AppConfig *Config

// LoadConfig 加载配置
func LoadConfig() {
	// 加载 .env 文件
	if err := godotenv.Load(); err != nil {
		log.Printf("未找到 .env 文件或读取失败: %v (将使用默认配置或环境变量)", err)
	}

	AppConfig = &Config{
		Host: getEnv("APP_HOST", "0.0.0.0"),
		Port: getEnv("APP_PORT", "8080"),
		Mode: getEnv("APP_MODE", "debug"), // debug 或 release
		DB: DatabaseConfig{
			Host:     getEnv("DB_HOST", "localhost"),
			Port:     getEnv("DB_PORT", "3306"),
			User:     getEnv("DB_USER", "root"),
			Password: getEnv("DB_PASSWORD", "123456"),
			DBName:   getEnv("DB_NAME", "projectTest"),
		},
		AI: AIConfig{
			APIKey:  getEnv("AI_API_KEY", "sk-api-myj2RN7e200E4y-Buq28YZgXeCUMj7ZmtGkRiZ-m9sdc3rxPnOx6lP4ILt685kM71H0jrc07OFPF7jHMT-VEsxKKkRzNnal-pXQw5SrnVTSC7GH6QN62O_M"),
			BaseURL: getEnv("AI_BASE_URL", "https://api.minimax.chat/v1"),
			Model:   getEnv("AI_MODEL", "abab6.5s-chat"),
		},
	}

	log.Println("配置加载成功")
}

// getEnv 获取环境变量，如果不存在则返回默认值
func getEnv(key, defaultValue string) string {
	value := os.Getenv(key)
	if value == "" {
		return defaultValue
	}
	return value
}
