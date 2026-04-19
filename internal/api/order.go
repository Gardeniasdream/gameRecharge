package api

import (
	"context"
	"fmt"
	"gameRecharge/config"
	"gameRecharge/internal/model"
	"gameRecharge/pkg/util"
	"time"

	"github.com/gin-gonic/gin"
)

func CreateOrder(c *gin.Context) {
	var req struct {
		UserId    uint64 `json:"user_id" binding:"required"`
		RoleId    string `json:"role_id" binding:"required"`
		ServerId  int    `json:"server_id" binding:"required"`
		ProductId string `json:"product_id" binding:"required"`
		PayType   int8   `json:"pay_type" binding:"required"` // 1微信 2支付宝
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(200, gin.H{"code": 1, "msg": "参数错误"})
		return
	}

	// 分布式锁：防止重复下单
	lockKey := fmt.Sprintf("lock:order:%d:%s", req.UserId, req.RoleId)
	ok, _ := config.Redis.SetNX(context.Background(), lockKey, 1, 10*time.Second).Result()
	if !ok {
		c.JSON(200, gin.H{"code": 1, "msg": "请勿重复提交"})
		return
	}
	defer config.Redis.Del(context.Background(), lockKey)

	// 生成订单信息
	orderNo := util.GenerateOrderNo()
	price, diamond, productName := util.GetProductInfo(req.ProductId)

	order := model.GameOrder{
		OrderNo:     orderNo,
		UserID:      req.UserId,
		RoleID:      req.RoleId,
		ServerID:    req.ServerId,
		ProductID:   req.ProductId,
		ProductName: productName,
		Price:       price,
		Diamond:     diamond,
		OrderStatus: 0,
		PayType:     req.PayType,
		ExpireTime:  time.Now().Add(30 * time.Minute),
		CreateTime:  time.Now(),
	}

	// 入库
	config.PG.Create(&order)

	// 返回支付信息
	c.JSON(200, gin.H{
		"code":     0,
		"msg":      "下单成功",
		"order_no": orderNo,
		"pay_url":  util.MockPayUrl(orderNo, price),
	})
}
