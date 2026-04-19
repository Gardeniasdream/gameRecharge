package util

import (
	"fmt"
	"time"
)

// GenerateOrderNo 生成唯一订单号
func GenerateOrderNo() string {
	return fmt.Sprintf("GM%d%d", time.Now().UnixMilli(), time.Now().Unix()%1000)
}
