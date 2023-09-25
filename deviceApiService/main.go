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
	pgHost      = "localhost"
	pgPort      = "15432"
	pgDB        = "devices"
	servicePort = ":8082"
)

func main() {

	//ctx := context.Background()

	router := gin.New()
	base := router.Group("/api/v1/device")
	config := base.Group("configProfile")

	dbconn := db.GetDb(pgHost, pgPort, pgDB)
	db.CreateProfileTable(dbconn)

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
