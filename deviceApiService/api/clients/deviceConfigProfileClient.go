package clients

import (
	"context"
	"net/http"

	"gorm.io/gorm"
	"proj.com/apisvc/db"
	"proj.com/apisvc/db/models"
)

type ConfigurationsClient interface {
	ReadProfile(context.Context, string) (models.DeviceConfigProfile, error)
	ListProfiles(context.Context) ([]models.DeviceConfigProfile, error)
	DeleteProfile(context.Context, string) error
	ApplyProfile(context.Context, string, string) error
	UploadProfileData(context.Context, []models.DeviceConfigProfile) error
}

type ConfigurationsClientImpl struct {
	HTTPClient http.Client
	PgDB       *gorm.DB
}

// NewConfigurationsClientImpl returns an impl of a configurations handler
func NewConfigurationsClientImpl(httpClient http.Client, pgDb *gorm.DB) ConfigurationsClient {
	return &ConfigurationsClientImpl{
		HTTPClient: httpClient,
		PgDB:       pgDb,
	}
}

func (impl *ConfigurationsClientImpl) ReadProfile(ctx context.Context, profileId string) (models.DeviceConfigProfile, error) {
	profile, err := db.ReadProfile(impl.PgDB, profileId)
	if err != nil {
		return models.DeviceConfigProfile{}, err
	}
	return profile, nil
}

func (impl *ConfigurationsClientImpl) ListProfiles(context.Context) ([]models.DeviceConfigProfile, error) {
	profiles, err := db.ReadAllProfiles(impl.PgDB)
	if err != nil {
		return []models.DeviceConfigProfile{}, err
	}
	return profiles, nil
}

func (impl *ConfigurationsClientImpl) DeleteProfile(ctx context.Context, profileId string) error {
	err := db.DeleteProfile(impl.PgDB, profileId)
	if err != nil {
		return err
	}
	return nil
}

func (impl *ConfigurationsClientImpl) UploadProfileData(ctx context.Context, profiles []models.DeviceConfigProfile) error {
	err := db.InsertProfile(impl.PgDB, profiles)
	if err != nil {
		return err
	}
	return nil
}

func (impl *ConfigurationsClientImpl) ApplyProfile(context.Context, string, string) error {
	return nil
}
