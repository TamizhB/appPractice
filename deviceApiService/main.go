package main

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"proj.com/apisvc/api/clients"
	"proj.com/apisvc/api/handlers"
	"proj.com/apisvc/db"
)

const (
	// need to push to env variables
	pgHost      = "postgres-service.default.svc.cluster.local"
	pgPort      = "15432"
	pgDB        = "devices"
	servicePort = ":8082"
)

func main() {

	//ctx := context.Background()

	router := gin.New()
	base := router.Group("/api/v1/device")

	dbconn := db.GetDb(pgHost, pgPort, pgDB)
	db.CreateProfileTable(dbconn)
	db.CreateDeviceTable(dbconn)

	deviceApiClient := clients.NewDeviceClientImpl(http.Client{}, dbconn)
	deviceApis := handlers.NewDeviceHandler(deviceApiClient)
	base.GET("/all", deviceApis.ListDevices)
	base.GET("/:device_name", deviceApis.GetDeviceConfig)
	base.DELETE("/:device_name", deviceApis.RemoveDevice)
	base.POST("/add", deviceApis.AddDevice)

	config := base.Group("configProfile")

	configApiClient := clients.NewConfigurationsClientImpl(http.Client{}, dbconn)
	configApis := handlers.NewConfigurationsHandler(configApiClient)

	config.GET("/:profile_name", configApis.ReadProfile)
	config.GET("/all", configApis.ListProfiles)
	config.DELETE("/:profile_name", configApis.DeleteProfile)
	config.POST("/upload", configApis.UploadProfileData)
	config.POST("/apply/:profile_name", configApis.ApplyProfile)

	server := &http.Server{Addr: servicePort, Handler: router}

	if err := server.ListenAndServe(); err != nil {
		panic(err)
	}
}
