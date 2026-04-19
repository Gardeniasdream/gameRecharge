package mq

import (
	"context"
	"encoding/json"
	"gameRecharge/config"
	"log"

	"github.com/apache/rocketmq-client-go/v2"
	"github.com/apache/rocketmq-client-go/v2/primitive"
	"github.com/apache/rocketmq-client-go/v2/producer"
)

// 全局生产者（单例）
var payProducer rocketmq.Producer

// PayMessage 支付消息结构体
type PayMessage struct {
	OrderNo    string `json:"order_no"`    // 订单号
	TradeNo    string `json:"trade_no"`    // 支付平台交易号
	PayChannel int8   `json:"pay_channel"` // 支付渠道 1微信 2支付宝
}

// InitRocketMQProducer 初始化生产者（在 main 中启动一次）
func InitRocketMQProducer() {
	namesrv := config.GetEnv("ROCKETMQ_ADDR")
	group := config.GetEnv("ROCKETMQ_GROUP")

	// 创建生产者
	p, err := rocketmq.NewProducer(
		producer.WithNameServer([]string{namesrv}),
		producer.WithGroupName(group),
		producer.WithRetry(3), // 失败重试3次
	)
	if err != nil {
		log.Fatal("RocketMQ 生产者初始化失败:", err)
	}

	// 启动生产者
	err = p.Start()
	if err != nil {
		log.Fatal("RocketMQ 生产者启动失败:", err)
	}

	payProducer = p
	log.Println("✅ RocketMQ 生产者启动成功")
}

// SendPayNotifyMessage 发送支付成功消息（给回调用）
func SendPayNotifyMessage(orderNo, tradeNo string, payChannel int8) error {
	topic := config.GetEnv("ROCKETMQ_TOPIC")

	// 构造消息体
	msgBody := PayMessage{
		OrderNo:    orderNo,
		TradeNo:    tradeNo,
		PayChannel: payChannel,
	}

	body, err := json.Marshal(msgBody)
	if err != nil {
		log.Println("消息序列化失败:", err)
		return err
	}

	// 构造RocketMQ消息
	msg := primitive.NewMessage(topic, body)

	// 发送消息
	_, err = payProducer.SendSync(context.Background(), msg)
	if err != nil {
		log.Println("发送支付消息失败:", err, "订单号:", orderNo)
		return err
	}

	log.Println("✅ 发送支付消息成功 订单号:", orderNo, " 交易号:", tradeNo)
	return nil
}
