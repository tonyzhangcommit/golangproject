package models

import "strconv"

// 用户表
type User struct {
	ID
	Name               string       `json:"name" gorm:"not null;column:name;comment:用户名"`
	Telnumber          string       `json:"telnumber" gorm:"column:telnumber;comment:电话"`
	Password           string       `json:"-" gorm:"column:password;comment:密码"`
	Status             bool         `json:"status" gorm:"column:status;default:true;comment:状态"`
	Roles              []Role       `gorm:"many2many:user_roles;"`
	Permissions        []Permission `gorm:"many2many:user_permission;"`
	IdentificationCode string       `json:"identificationcode" gorm:"unique;column:identificationcode;default:NULL;comment:推广码"`
	ManagerID          uint
	UserList           []User       `json:"userlist" gorm:"foreignkey:ManagerID"`
	OrderList          []Order      `json:"orderlist"`
	InComeList         []InComeInfo `json:"inComelist"`
	Timestamps
}

func (user User) GetUid() string {
	return strconv.Itoa(int(user.ID.ID))
}

// 角色表
type Role struct {
	ID
	Name string `json:"name" gorm:"not null;unique;column:name;comment:角色名"`
	Timestamps
}

// 权限表
type Permission struct {
	ID
	Name string `json:"name" gorm:"not null;column:name;comment:权限名"`
	Timestamps
}
