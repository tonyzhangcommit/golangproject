package models

// 店铺表
type Shop struct {
	ID
	UserID     uint
	Name       string     `json:"name" gorm:"not null;column:name;comment:店铺名"`
	Address    string     `json:"address" gorm:"not null;column:address;comment:店铺地址"`
	Tables     []Table    `gorm:"constraint:OnUpdate:CASCADE,OnDelete:SET NULL;"`
	Catagories []Catagory `gorm:"constraint:OnUpdate:CASCADE,OnDelete:SET NULL;"`
	Timestamps
}

// 店铺桌子表信息（扫码点餐需要配置）
type Table struct {
	ID
	Name     string `json:"name" gorm:"not null;column:name;comment:桌名"`
	TableNum uint32 `json:"tablenum" gorm:"column:tablenum;comment:桌号"`
	QRCode   string `json:"qrcode" gorm:"column:qrcode;comment:二维码"`
	Status   string `json:"status" gorm:"column:status;comment:状态"`
	ShopID   uint
	Timestamps
}

// 菜品分类
type Catagory struct {
	ID
	Name     string    `json:"name" gorm:"not null;column:name;comment:分类名"`
	Cuisines []Cuisine `gorm:"many2many:Catagory_Cuisine;"`
	ShopID   uint
	Timestamps
}

type Cuisine struct {
	ID
	Name        string     `json:"name" gorm:"not null;column:name;comment:菜名"`
	Price       float32    `json:"price" gorm:"not null;column:price;comment:价格"`
	ActualPrice float32    `json:"actualprice" gorm:"not null;column:actualprice;comment:实际价格"`
	Sales       uint32     `json:"sales" gorm:"column:sales;comment:销量"`                 // 统计月销量
	Repeat      uint32     `json:"repeat" gorm:"column:repeat;comment:回头客"`             // 统计回头客
	Peculiarity string     `json:"peculiarity" gorm:"column:peculiarity;comment:特色简介"` // 统计回头客
	Catagories  []Catagory `gorm:"many2many:Catagory_Cuisine;"`
	Orders  []Order `gorm:"many2many:Order_Cuisine;"`
	Timestamps
}
