package models

// 用户表
type User struct {
	ID
	Name      string `json:"name" gorm:"not null;column:name;comment:用户名"`
	Telnumber string `json:"telnumber" gorm:"column:telnumber;comment:电话"`
	Password  string `json:"password" gorm:"column:password;comment:密码"`
	Status    bool   `json:"status" gorm:"column:status;default:true;comment:状态"`
	Timestamps
	Roles   []Role    `gorm:"many2many:user_roles;"`
	Coupons []*Coupon `gorm:"many2many:user_coupons;"`
	Shops   []Shop    `gorm:"constraint:OnUpdate:CASCADE,OnDelete:SET NULL;"`
	Orders  []Order   `gorm:"constraint:OnUpdate:CASCADE,OnDelete:SET NULL;"`
}

// 角色表
type Role struct {
	ID
	Name        string       `json:"name" gorm:"not null;column:name;comment:角色名"`
	Users       []User       `gorm:"many2many:user_roles;"`
	Permissions []Permission `gorm:"many2many:Permission_roles;"`
	Timestamps
}

// 权限表
type Permission struct {
	ID
	Name  string `json:"name" gorm:"not null;column:name;comment:权限名"`
	Roles []Role `gorm:"many2many:Permission_roles;"`
	Timestamps
}
