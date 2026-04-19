package model

import "time"

// GameOrder 游戏充值订单表
type GameOrder struct {
	ID          uint64    `gorm:"primaryKey;column:id;type:bigserial"`
	OrderNo     string    `gorm:"column:order_no;type:varchar(64);uniqueIndex;not null"`
	UserID      uint64    `gorm:"column:user_id;type:bigint;not null"`
	RoleID      string    `gorm:"column:role_id;type:varchar(64);not null"`
	ServerID    int       `gorm:"column:server_id;type:int;not null"`
	ProductID   string    `gorm:"column:product_id;type:varchar(32);not null"`
	ProductName string    `gorm:"column:product_name;type:varchar(64);not null"`
	Price       float64   `gorm:"column:price;type:decimal(10,2);not null"`
	Diamond     int       `gorm:"column:diamond;type:int;not null"`
	OrderStatus int8      `gorm:"column:order_status;type:int2;not null;default:0"` // 0待支付 1已支付 2取消
	PayType     int8      `gorm:"column:pay_type;type:int2;not null"`               // 1微信 2支付宝
	PayTime     time.Time `gorm:"column:pay_time;type:timestamp"`
	ExpireTime  time.Time `gorm:"column:expire_time;type:timestamp;not null"`
	CreateTime  time.Time `gorm:"column:create_time;type:timestamp;default:now()"`
}

func (GameOrder) TableName() string {
	return "game_order"
}

// PayLog 支付流水表
type PayLog struct {
	ID            uint64    `gorm:"primaryKey;column:id;type:bigserial"`
	OrderNo       string    `gorm:"column:order_no;type:varchar(64);not null"`
	TransactionID string    `gorm:"column:transaction_id;type:varchar(64);default:''"`
	PayChannel    int8      `gorm:"column:pay_channel;type:int2;not null"`
	Amount        float64   `gorm:"column:amount;type:decimal(10,2);not null"`
	NotifyRaw     string    `gorm:"column:notify_raw;type:text"`
	CreateTime    time.Time `gorm:"column:create_time;type:timestamp;default:now()"`
}

func (PayLog) TableName() string {
	return "pay_log"
}

// RewardFlow 钻石发放流水表
type RewardFlow struct {
	ID         uint64    `gorm:"primaryKey;column:id;type:bigserial"`
	OrderNo    string    `gorm:"column:order_no;type:varchar(64);not null"`
	RoleID     string    `gorm:"column:role_id;type:varchar(64);not null"`
	ServerID   int       `gorm:"column:server_id;type:int;not null"`
	Diamond    int       `gorm:"column:diamond;type:int;not null"`
	Status     int8      `gorm:"column:status;type:int2;default:0"` // 0待发放 1已发放
	CreateTime time.Time `gorm:"column:create_time;type:timestamp;default:now()"`
}

func (RewardFlow) TableName() string {
	return "reward_flow"
}

// Idempotent 幂等防重表
type Idempotent struct {
	UniqueKey string    `gorm:"primaryKey;column:unique_key;type:varchar(64)"`
	CreatedAt time.Time `gorm:"column:created_at;type:timestamp;default:now()"`
}

func (Idempotent) TableName() string {
	return "idempotent"
}
