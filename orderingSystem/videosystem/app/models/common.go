package models

import (
	"time"

	"gorm.io/gorm"
)

type ID struct {
	ID uint `json:"id" gorm:"primaryKey"`
}

// 创建、更新时间
type Timestamps struct {
	CreatedAt time.Time `json:"-"`
	UpdatedAt time.Time `json:"-"`
}

// 软删除
type SoftDeletes struct {
	DeletedAt gorm.DeletedAt `json:"deleted_at" gorm:"index"`
}
