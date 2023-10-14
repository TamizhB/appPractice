package models

import "gorm.io/gorm"

type Device struct {
	gorm.Model
	ID             string `gorm:"primary_key"`
	Address        string
	Port           string
	AppliedProfile DeviceConfigProfile `gorm:"foreignKey:ID"`
	RecentConfig   string
}
