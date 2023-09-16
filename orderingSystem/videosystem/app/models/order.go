package models

import "time"

// 产品表 周，月，季，年会员
type Products struct {
	ID
	Name  string  `json:"name" gorm:"not null;column:name;comment:会员名称"`
	Price float32 `json:"price" gorm:"not null;column:price;comment:会员价格"`
	Timestamps
}

// 会员表  订单类型，订单状态，用户，产品(只有一对一关系)
type Order struct {
	ID
	UserID        uint
	OrderCategory string  `json:"ordercategory" gorm:"not null;column:ordercategory;comment:订单类型"`
	OrderStatus   string  `json:"orderstatus" gorm:"not null;column:orderstatus;comment:订单状态"`
	Product       uint    `json:"product" gorm:"not null;column:product;comment:下单物品"`
	Price         float32 `json:"price" gorm:"not null;column:price;comment:支付金额"`
	OrderInfo     OrderInfo
}

type OrderInfo struct {
	ID
	OrderID   uint
	PayType   string    `json:"paytype" gorm:"not null;column:paytype;comment:支付方式"`
	CreatedAt time.Time `json:"datetime" gorm:"not null;column:datetime;comment:支付日期"`
}

// 收益表 收益类型，代理收入/直接收入
type InComeInfo struct {
	ID
	UserID     uint
	InComeType string
	InComeNum  float64
	CreatedAt  time.Time `json:"createT" gorm:"not null;column:datetime;comment:支付日期"`
}
