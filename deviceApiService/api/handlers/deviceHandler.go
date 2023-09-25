package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type DeviceHandler struct {
	HTTPClient http.Client
}

// NewDeviceHandler returns an impl of a device handler
func NewDeviceHandler(httpClient http.Client) DeviceHandler {
	return DeviceHandler{
		HTTPClient: httpClient,
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
