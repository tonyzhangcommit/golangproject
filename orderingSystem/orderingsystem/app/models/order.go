package models

import "time"

// 暂定一种满减优惠券
type Coupon struct {
	ID
	Code      string    `json:"code" gorm:"not null;unique;column:code;comment:券码"`
	Condition string    `json:"condition" gorm:"not null;column:condition;comment:名称"`
	MiniPrice float64   `json:"miniprice" gorm:"not null;column:miniprice;comment:使用条件"`
	ExpiredAt time.Time `json:"expiredtime" gorm:"not null;column:expiredtime;comment:过期时间"`
	OrderID   uint
	Timestamps
}

// 订单表
type Order struct {
	ID
	UserID      uint
	TableID     uint
	Cuisines    []Cuisine `gorm:"many2many:Order_Cuisine"`
	Price       float32   `json:"price" gorm:"not null;column:price;comment:价格"`
	ActualPrice float32   `json:"actualprice" gorm:"not null;column:actualprice;comment:实际价格"`
	Coupons     []*Coupon `gorm:"constraint:OnUpdate:CASCADE,OnDelete:SET NULL;"`
	Status      string    `json:"status" gorm:"column:status;comment:状态"`
	PayWay      string    `json:"payway" gorm:"column:payway;comment:支付方式"`
	Timestamps
}
