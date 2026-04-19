package mq

import (
	"context"
	"encoding/json"
	"gameRecharge/config"
	"gameRecharge/internal/model"
	"gameRecharge/pkg/gamesrv"
	"log"
	"time"

	"github.com/apache/rocketmq-client-go/v2"
	"github.com/apache/rocketmq-client-go/v2/consumer"
	"github.com/apache/rocketmq-client-go/v2/primitive"
	"gorm.io/gorm"
)

// 启动支付消息消费者
func ConsumePayNotify() {
	namesrv := config.GetEnv("ROCKETMQ_ADDR")
	group := config.GetEnv("ROCKETMQ_GROUP")
	topic := config.GetEnv("ROCKETMQ_TOPIC")

	// 创建消费者
	c, err := rocketmq.NewPushConsumer(
		consumer.WithNameServer([]string{namesrv}),
		consumer.WithGroupName(group),
		consumer.WithConsumeFromWhere(consumer.ConsumeFromLastOffset),
	)
	if err != nil {
		log.Fatal("创建RocketMQ消费者失败:", err)
	}

	// 订阅topic + 消息处理
	err = c.Subscribe(topic, consumer.MessageSelector{}, func(ctx context.Context,
		msgs ...*primitive.MessageExt) (consumer.ConsumeResult, error) {

		for _, msg := range msgs {
			log.Println("收到支付消息:", string(msg.Body))

			// 1. 解析消息
			var payMsg PayMessage
			if err := json.Unmarshal(msg.Body, &payMsg); err != nil {
				log.Println("消息解析失败", err)
				return consumer.ConsumeSuccess, nil
			}

			// 2. 执行核心逻辑（幂等+事务+发钻石）
			HandlePaySuccess(payMsg)
		}

		return consumer.ConsumeSuccess, nil
	})

	if err != nil {
		log.Fatal("订阅Topic失败:", err)
	}

	// 启动消费者
	_ = c.Start()
	defer c.Shutdown()

	log.Println("✅ RocketMQ 支付消费者启动成功，监听中...")

	// 阻塞
	select {}
}

// HandlePaySuccess 处理支付成功（核心：发钻石）
func HandlePaySuccess(msg PayMessage) {
	orderNo := msg.OrderNo
	tradeNo := msg.TradeNo
	payChannel := msg.PayChannel

	// ==========================================
	// 【1】幂等校验（绝对防重复）
	// ==========================================
	var count int64
	config.PG.Model(&model.Idempotent{}).
		Where("unique_key = ?", tradeNo).
		Count(&count)
	if count > 0 {
		log.Println("⚠️ 订单已处理，跳过重复消息:", orderNo)
		return
	}

	// ==========================================
	// 【2】事务：订单 + 支付流水 + 钻石流水 + 幂等（一起成功/失败）
	// ==========================================
	err := config.PG.Transaction(func(tx *gorm.DB) error {
		// --------------------------
		// A. 更新订单为已支付
		// --------------------------
		res := tx.Model(&model.GameOrder{}).
			Where("order_no = ? AND order_status = 0", orderNo).
			Updates(map[string]interface{}{
				"order_status": 1,
				"pay_time":     time.Now(),
			})
		if res.RowsAffected == 0 {
			log.Println("订单已支付或不存在:", orderNo)
			return nil
		}

		// --------------------------
		// B. 查询订单信息
		// --------------------------
		var order model.GameOrder
		if err := tx.Where("order_no = ?", orderNo).First(&order).Error; err != nil {
			return err
		}

		// --------------------------
		// C. 写入【支付流水】pay_log
		// --------------------------
		payLog := model.PayLog{
			OrderNo:       orderNo,
			TransactionID: tradeNo,
			PayChannel:    payChannel,
			Amount:        order.Price,
			NotifyRaw:     "",
			CreateTime:    time.Now(),
		}
		if err := tx.Create(&payLog).Error; err != nil {
			return err
		}

		// --------------------------
		// D. 写入【钻石发放流水】
		// --------------------------
		reward := model.RewardFlow{
			OrderNo:    orderNo,
			RoleID:     order.RoleID,
			ServerID:   order.ServerID,
			Diamond:    order.Diamond,
			Status:     0,
			CreateTime: time.Now(),
		}
		if err := tx.Create(&reward).Error; err != nil {
			return err
		}

		// --------------------------
		// E. 写入幂等记录
		// --------------------------
		idem := model.Idempotent{
			UniqueKey: tradeNo,
			CreatedAt: time.Now(),
		}
		return tx.Create(&idem).Error
	})

	if err != nil {
		log.Println("事务执行失败:", err)
		return
	}

	// ==========================================
	// 【3】调用游戏服发钻石
	// ==========================================
	var order model.GameOrder
	config.PG.Where("order_no = ?", orderNo).First(&order)
	gamesrv.SendDiamond(order.ServerID, order.RoleID, order.Diamond)

	// 更新发放状态
	config.PG.Model(&model.RewardFlow{}).
		Where("order_no = ?", orderNo).
		Update("status", 1)

	log.Println("✅ 订单处理完成，钻石已发放:", orderNo)
}
