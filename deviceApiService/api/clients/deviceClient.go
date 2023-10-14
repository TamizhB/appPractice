package clients

import (
	"context"
	"net/http"

	"gorm.io/gorm"
	"proj.com/apisvc/db"
	"proj.com/apisvc/db/models"
)

type DeviceClient interface {
	AddDevice(ctx context.Context, device models.Device) error
	RemoveDevice(ctx context.Context, deviceId string) error
	ListDevices(context.Context) ([]models.Device, error)
	GetDeviceConfig(ctx context.Context, deviceId string) (models.Device, error)
}

type DeviceClientImpl struct {
	HTTPClient http.Client
	PgDB       *gorm.DB
}

// NewConfigurationsClientImpl returns an impl of a configurations handler
func NewDeviceClientImpl(httpClient http.Client, pgDb *gorm.DB) DeviceClient {
	return &DeviceClientImpl{
		HTTPClient: httpClient,
		PgDB:       pgDb,
	}
}

func (impl *DeviceClientImpl) GetDeviceConfig(ctx context.Context, deviceId string) (models.Device, error) {
	device, err := db.GetDevice(impl.PgDB, deviceId)
	if err != nil {
		return models.Device{}, err
	}
	return device, nil
}

func (impl *DeviceClientImpl) ListDevices(context.Context) ([]models.Device, error) {
	devices, err := db.ReadAllDevices(impl.PgDB)
	if err != nil {
		return []models.Device{}, err
	}
	return devices, nil
}

func (impl *DeviceClientImpl) RemoveDevice(ctx context.Context, deviceId string) error {
	err := db.DeleteDevice(impl.PgDB, deviceId)
	if err != nil {
		return err
	}
	return nil
}

func (impl *DeviceClientImpl) AddDevice(ctx context.Context, device models.Device) error {
	err := db.InsertDevice(impl.PgDB, device)
	if err != nil {
		return err
	}
	return nil
}
