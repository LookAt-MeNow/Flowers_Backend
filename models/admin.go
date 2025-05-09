package models

import "gorm.io/gorm"

type Admin struct {
	gorm.Model
	Username string `gorm:"uniqueIndex;size:50;not null"`
	Password string `gorm:"size:255;not null"`
	Role     string `gorm:"size:20;default:'admin'"`
}