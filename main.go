package main

import (
	"gameRecharge/config"
)

func main() {
	// 初始化配置、PG、Redis
	config.Init()

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
}
