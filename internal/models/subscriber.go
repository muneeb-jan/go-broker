package models

import "gorm.io/gorm"

type Subscriber struct {
	gorm.Model
	ID       string `gorm:"uniqueIndex"`
	Topic    string
	Listener string
}
