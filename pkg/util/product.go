package util

import "fmt"

// GetProductInfo 获取充值档位：金额、钻石、名称
func GetProductInfo(productId string) (float64, int, string) {
	switch productId {
	case "1":
		return 6.0, 60, "6元=60钻石"
	case "2":
		return 30.0, 330, "30元=330钻石"
	case "3":
		return 68.0, 800, "68元=800钻石"
	default:
		return 0, 0, "未知档位"
	}
}

// MockPayUrl 模拟支付链接
func MockPayUrl(orderNo string, price float64) string {
	return fmt.Sprintf("https://pay.xxx.com?order_no=%s&amount=%.2f", orderNo, price)
}
