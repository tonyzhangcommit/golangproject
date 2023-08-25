package models

// 视频简介
type Video struct {
	ID
	Name       string      `json:"name" gorm:"not null;column:name;comment:剧名"`
	Cover      string      `json:"cover" gorm:"not null;column:cover;comment:封面"`
	Intro      string      `json:"intro" gorm:"not null;column:intro;comment:简介"`
	Categories []Category  `gorm:"many2many:video_category;"`
	Videolist  []VideoInfo `gorm:"constraint:OnUpdate:CASCADE,OnDelete:SET NULL;"`
	Timestamps
}

// 视频详情
type VideoInfo struct {
	ID
	VideoID  uint
	Episodes int    `json:"episodes" gorm:"column:episodes;comment:剧集"`
	Url      string `json:"url" gorm:"not null;column:url;comment:视频地址"`
	Intro    string `json:"intro" gorm:"default:'';column:intro;comment:简介"`
	Timestamps
}

// 分类
type Category struct {
	ID
	Name         string `json:"name" gorm:"unique;not null;column:name;comment:名称"`
	Intro        string `json:"intro" gorm:"not null;column:intro;comment:简介;defaule:''"`
	CategoryID   uint
	Categorylist []Category `gorm:"foreignkey:CategoryID"`
	Timestamps
}
