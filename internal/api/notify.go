package api

import (
	"gameRecharge/internal/mq"
	"net/http"

	"github.com/gin-gonic/gin"
)

// PayNotify 支付回调接口
func PayNotify(c *gin.Context) {
	// 1. 接收支付平台参数
	orderNo := c.PostForm("out_trade_no")
	tradeNo := c.PostForm("transaction_id")
	payChannel := 1 // 1微信 2支付宝（根据实际渠道改）

	// 2. 发送 MQ 消息（异步处理发钻石）
	_ = mq.SendPayNotifyMessage(orderNo, tradeNo, int8(payChannel))

	// 3. 必须返回 success 给微信/支付宝
	c.JSON(http.StatusOK, gin.H{"code": 0, "msg": "success"})
}
