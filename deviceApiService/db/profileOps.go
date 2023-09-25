package db

import (
	"gorm.io/gorm"
	"proj.com/apisvc/db/models"
)

func CreateProfileTable(db *gorm.DB) error {
	return db.AutoMigrate(&models.DeviceConfigProfile{})
}

func InsertProfile(db *gorm.DB, profiles []models.DeviceConfigProfile) error {
	result := db.CreateInBatches(&profiles, 10)
	return result.Error
}

func ReadProfile(db *gorm.DB, profileId string) (models.DeviceConfigProfile, error) {
	var profiles []models.DeviceConfigProfile
	result := db.Where("ID = ?", profileId).Find(&profiles)
	if result.Error != nil || len(profiles) == 0 {
		return models.DeviceConfigProfile{}, result.Error
	}
	return profiles[0], nil
}

func ReadAllProfiles(db *gorm.DB) ([]models.DeviceConfigProfile, error) {
	var profiles []models.DeviceConfigProfile
	result := db.Find(&profiles)
	return profiles, result.Error
}

func DeleteProfile(db *gorm.DB, profileId string) error {
	result := db.Where("ID = ?", profileId).Delete(&models.DeviceConfigProfile{})
	return result.Error
}

func UpdateProfile(db *gorm.DB, profile models.DeviceConfigProfile) error {
	result := db.Where("ID = ?", profile.ID).Save(&models.DeviceConfigProfile{})
	return result.Error
}
