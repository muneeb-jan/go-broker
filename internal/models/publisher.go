package models

import "gorm.io/gorm"

type Publisher struct {
	gorm.Model
	ID string `gorm:"uniqueIndex"`
}
