package request

type Product struct {
	PID   uint    `form:"productid" json:"productid"`
	Name  string  `form:"name" json:"name"`
	Price float32 `form:"price" json:"price"`
}

type Delp struct {
	Pid uint `form:"pid" json:"pid"`
}

type CreateOrder struct {
	UserID        uint    `form:"userid" json:"userid"`                                  // 用户ID
	OrderCategory string  `form:"ordercategory" json:"ordercategory" binding:"required"` // 订单类型：menbershiptype  videotype
	Product       uint    `form:"product" json:"product" binding:"required"`             // 视频ID或者产品ID
	Price         float32 `form:"price" json:"price" binding:"required"`                 // 价格
}

type Pay struct {
	UserID  uint   `form:"userid" json:"userid"`
	OrderID uint   `form:"orderid" json:"orderid" binding:"required"`
	PayType string `form:"paytype" json:"paytype" binding:"required"`
}

type ChangeOrderStatus struct {
	OrderID uint `form:"orderid" json:"orderid" binding:"required"`
	UserID  uint `form:"userid" json:"userid" binding:"required"`
}
