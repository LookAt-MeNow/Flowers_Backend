package models

import (
	"gorm.io/gorm"
)

type Merchant struct {
	gorm.Model
	Username  string    `gorm:"uniqueIndex;size:50;not null"`
	Password  string    `gorm:"size:255;not null"`
	ShopName  string    `gorm:"size:100;not null"`
	Email     string    `gorm:"uniqueIndex;size:100;not null"`
	Phone     string    `gorm:"size:20;not null"`
	Address   string    `gorm:"size:255"`
	Status    int       `gorm:"default:1"` // 1-正常, 0-禁用
}