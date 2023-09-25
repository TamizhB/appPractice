package models

import "gorm.io/gorm"

type DeviceConfigProfile struct {
	gorm.Model
	ID   string `gorm:"primary_key"`
	Data string
}
