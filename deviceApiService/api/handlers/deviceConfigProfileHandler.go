package handlers

import (
	"fmt"
	"io"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
	"proj.com/apisvc/api/clients"
	"proj.com/apisvc/messageHandlers"
	"proj.com/apisvc/parsers"
)

type ConfigurationsHandler struct {
	ConfClient clients.ConfigurationsClient
}

// NewConfigurationsHandler returns an impl of a configurations handler
func NewConfigurationsHandler(confClient clients.ConfigurationsClient) ConfigurationsHandler {
	return ConfigurationsHandler{
		ConfClient: confClient,
	}
}

func (h *ConfigurationsHandler) ReadProfile(ctx *gin.Context) {
	profileName := ctx.Param("profile_name")
	profile, err := h.ConfClient.ReadProfile(ctx, profileName)
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusInternalServerError, err.Error())
		return
	}
	ctx.JSON(http.StatusOK, profile)
}

func (h *ConfigurationsHandler) DeleteProfile(ctx *gin.Context) {
	profileName := ctx.Param("profile_name")
	err := h.ConfClient.DeleteProfile(ctx, profileName)
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusInternalServerError, err.Error())
		return
	}
	ctx.JSON(http.StatusOK, "Profile "+profileName+" deleted successfully")
}

func (h *ConfigurationsHandler) ListProfiles(ctx *gin.Context) {
	profiles, err := h.ConfClient.ListProfiles(ctx)
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusInternalServerError, err.Error())
		return
	}
	ctx.JSON(http.StatusOK, profiles)
}

func (h *ConfigurationsHandler) UploadProfileData(ctx *gin.Context) {
	file, _, err := ctx.Request.FormFile("file")
	if err != nil {
		ctx.AbortWithStatusJSON(404, fmt.Errorf("cannot get the file"))
		return
	}
	defer file.Close()

	tempFile, err := os.Create("temp.xlsx")
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusInternalServerError, err.Error())
		return
	}
	defer os.Remove("temp.xlsx")

	_, err = io.Copy(tempFile, file)
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusInternalServerError, err.Error())
		return
	}

	profileList, err := parsers.ParseProfileExcel("temp.xlsx", 0)
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusInternalServerError, err.Error())
		return
	}

	err = h.ConfClient.UploadProfileData(ctx, profileList)
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusInternalServerError, err.Error())
		return
	}

	ctx.JSON(http.StatusCreated, "Config profiles uploaded successfully")
}

func (h *ConfigurationsHandler) ApplyProfile(ctx *gin.Context) {
	fmt.Printf("TBD")
}

func (h *ConfigurationsHandler) ApplyConfig(ctx *gin.Context) {
	body, err := io.ReadAll(ctx.Request.Body)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Error reading request body"})
		return
	}

	jsonPayload := string(body)
	fmt.Println("Request:", jsonPayload)
	respStatus, respBody, err := messageHandlers.Send("configRequest", jsonPayload)
	fmt.Println("Status:", respStatus)
	fmt.Println("Message:", respBody)

	if err != nil {
		ctx.JSON(respStatus, err.Error())
		return
	}

	ctx.JSON(respStatus, respBody)
}

// Message struct definition
type Message struct {
	MessageType     int    `json:"MessageType"`
	ActualCommand   int    `json:"ActualCommand"`
	SenderID        uint32 `json:"SenderID"`
	ReceiverID      uint32 `json:"ReceiverID"`
	MsgID           uint32 `json:"MsgID"`
	SequenceNumber  uint32 `json:"SequenceNumber"`
	TimestampSender int64  `json:"TimestampSender"`
	Payload         string `json:"Payload"`
}
