package entities

import "github.com/jinzhu/gorm"

type (
	File struct {
		gorm.Model
		Status string `gorm:"size:16;default:'PENDING'"`
	}
)
