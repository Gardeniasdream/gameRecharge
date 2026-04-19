package main

import (
	"gameRecharge/config"
	"gameRecharge/internal/api"
	"gameRecharge/internal/mq"

	"github.com/gin-gonic/gin"
)

func main() {
	// 初始化配置、PG、Redis
	config.Init()

	mq.InitRocketMQProducer() // 初始化 RocketMQ 生产者

	// 启动 MQ 消费者
	go mq.ConsumePayNotify()

	// 初始化数据库后执行
	//err := config.PG.AutoMigrate(
	//	&model.GameOrder{},
	//	&model.PayLog{},
	//	&model.RewardFlow{},
	//	&model.Idempotent{},
	//)
	//if err != nil {
	//	log.Println("初始化数据库表失败")
	//	return
	//}

	// 启动 HTTP 服务
	r := gin.Default()

	// 接口路由
	r.POST("/recharge/create", api.CreateOrder)
	r.POST("/pay/notify", api.PayNotify)

	// 启动端口
	port := config.GetEnv("SERVER_PORT")
	r.Run(":" + port)

}
