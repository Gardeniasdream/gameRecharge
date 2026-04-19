package config

import (
	"os"
	"strconv"

	"github.com/joho/godotenv"
)

// LoadEnv 加载 .env 文件
func LoadEnv() {
	_ = godotenv.Load(".env")
}

// GetEnv 获取字符串配置
func GetEnv(key string) string {
	return os.Getenv(key)
}

// GetEnvInt 获取数字配置
func GetEnvInt(key string) int {
	val, err := strconv.Atoi(os.Getenv(key))
	if err != nil {
		return 0
	}
	return val
}
