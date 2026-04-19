package config

import (
	"context"
	"fmt"
	"log"

	"github.com/go-redis/redis/v8"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

// 全局数据库、Redis 客户端
var (
	PG    *gorm.DB
	Redis *redis.Client
)

// Init 初始化所有中间件
func Init() {
	LoadEnv()
	InitPostgreSQL()
	InitRedis()
}

// InitPostgreSQL 初始化 PostgreSQL
func InitPostgreSQL() {
	dsn := fmt.Sprintf(
		"host=%s user=%s password=%s dbname=%s port=%d sslmode=%s TimeZone=Asia/Shanghai",
		GetEnv("PG_HOST"),
		GetEnv("PG_USER"),
		GetEnv("PG_PASSWORD"),
		GetEnv("PG_DATABASE"),
		GetEnvInt("PG_PORT"),
		GetEnv("PG_SSL_MODE"),
	)

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info),
	})
	if err != nil {
		log.Fatal("PostgreSQL 连接失败：", err)
	}

	PG = db
	log.Println("PostgreSQL 连接成功")
}

// InitRedis 初始化 Redis
func InitRedis() {
	client := redis.NewClient(&redis.Options{
		Addr:     GetEnv("REDIS_ADDR"),
		Password: GetEnv("REDIS_PASSWORD"),
		DB:       GetEnvInt("REDIS_DB"),
	})

	_, err := client.Ping(context.Background()).Result()
	if err != nil {
		log.Fatal("Redis 连接失败：", err)
	}

	Redis = client
	log.Println("Redis 连接成功")
}
