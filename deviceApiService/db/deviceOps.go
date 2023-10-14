package db

import (
	"gorm.io/gorm"
	"proj.com/apisvc/db/models"
)

func CreateDeviceTable(db *gorm.DB) error {
	return db.AutoMigrate(&models.Device{})
}

func InsertDevice(db *gorm.DB, device models.Device) error {
	result := db.Create(&device)
	return result.Error
}

func GetDevice(db *gorm.DB, deviceId string) (models.Device, error) {
	var devices []models.Device
	result := db.Where("ID = ?", deviceId).Find(&devices)
	if result.Error != nil || len(devices) == 0 {
		return models.Device{}, result.Error
	}
	return devices[0], nil
}

func ReadAllDevices(db *gorm.DB) ([]models.Device, error) {
	var devices []models.Device
	result := db.Find(&devices)
	return devices, result.Error
}

func DeleteDevice(db *gorm.DB, deviceId string) error {
	result := db.Where("ID = ?", deviceId).Delete(&models.Device{})
	return result.Error
}

func UpdateDevice(db *gorm.DB, device models.Device) error {
	result := db.Where("ID = ?", device.ID).Save(&models.Device{})
	return result.Error
}
