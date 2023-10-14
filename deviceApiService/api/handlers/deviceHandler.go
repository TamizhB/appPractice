package handlers

import (
	"github.com/gin-gonic/gin"
	"proj.com/apisvc/api/clients"
)

type DeviceHandler struct {
	DevClient clients.DeviceClient
}

// NewDeviceHandler returns an impl of a device handler
func NewDeviceHandler(devClient clients.DeviceClient) DeviceHandler {
	return DeviceHandler{
		DevClient: devClient,
	}
}

func (h *DeviceHandler) AddDevice(ctx *gin.Context) {

}

func (h *DeviceHandler) RemoveDevice(ctx *gin.Context) {

}

func (h *DeviceHandler) GetDevice(ctx *gin.Context) {

}

func (h *DeviceHandler) ListDevices(ctx *gin.Context) {

}

func (h *DeviceHandler) GetDeviceConfig(ctx *gin.Context) {

}
