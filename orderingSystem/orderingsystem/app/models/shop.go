package models

// 店铺表
type Shop struct {
	ID
	UserID     uint       `json:"userid" gorm:"column:user_id;comment:店主ID"`
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
	Cuisines []Cuisine `gorm:"constraint:OnUpdate:CASCADE,OnDelete:SET NULL;"`
	ShopID   uint      `gorm:"column:shop_id;comment:店铺ID"`
	Timestamps
}

type Image struct {
	ID
	CuisineID uint
	ImageUrl  string `json:"imageurl" gorm:"not null;column:imageurl;comment:图片URL"`
	Timestamps
}

// 菜品详情
type Cuisine struct {
	ID
	CatagoryID  uint
	Name        string  `json:"name" gorm:"not null;column:name;comment:菜名"`
	Images      []Image `gorm:"constraint:OnUpdate:CASCADE,OnDelete:SET NULL;"`
	Price       float32 `json:"price" gorm:"not null;column:price;comment:价格"`
	Sales       uint64  `json:"sales" gorm:"column:sales;comment:销量"`               // 统计月销量
	Repeat      uint64  `json:"repeat" gorm:"column:repeat;comment:回头客"`            // 统计回头客
	Peculiarity string  `json:"peculiarity" gorm:"column:peculiarity;comment:特色简介"`
	Timestamps
}
