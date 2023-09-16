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
	UserID        uint   `form:"userid" json:"userid"`
	OrderCategory string `form:"ordercategory" json:"ordercategory" binding:"required"`
	Product       uint
	Price         float32 `form:"price" json:"price" binding:"required"`
}
