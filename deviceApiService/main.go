package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	tracesdk "go.opentelemetry.io/otel/sdk/trace"
	"proj.com/apisvc/api/clients"
	"proj.com/apisvc/api/handlers"
	"proj.com/apisvc/db"
	"proj.com/apisvc/tracing"
)

const (
	// need to push to env variables
	pgHost      = "postgres-service.default.svc.cluster.local"
	pgPort      = "5432"
	pgDB        = "devices"
	servicePort = ":8082"
)

var traceProvider *tracesdk.TracerProvider

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
	config.POST("/apply", configApis.ApplyConfig)

	server := &http.Server{Addr: servicePort, Handler: router}

	if err := server.ListenAndServe(); err != nil {
		panic(err)
	}
}

func initTracing() {
	// init tracing
	jaegerEndpoint := fmt.Sprintf("%s%s", strings.TrimRight("http://simplest-collector.default.svc.cluster.local:14268", "/"), "/api/traces")
	var err error
	ctx := context.Background()
	traceProvider, err = tracing.TracerProvider(ctx, 1, jaegerEndpoint, "api-svc", "dev")
	if err != nil {
		log.Println("unable to start tracing")
	}
}
