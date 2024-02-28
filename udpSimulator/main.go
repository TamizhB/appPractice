package main

import (
	"context"
	"fmt"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"proj.com/udp/messageHandlers"
)

const (
	serverAddr   = "127.0.0.1"
	messageCount = 10
	servicePort  = ":3000"
)

func main() {
	go messageHandlers.Serve()

	router := gin.New()
	base := router.Group("/udpSvc")

	server := &http.Server{Addr: servicePort, Handler: router}
	msg := base.Group("msg")
	msg.POST("/send", messageHandlers.SendMessage)
	initTracing()
	if err := server.ListenAndServe(); err != nil {
		panic(err)
	}
}

func initTracing() {
	// init tracing
	jaegerEndpoint := fmt.Sprintf("%s%s", strings.TrimRight(envConfig.JaegerAPIEndpoint, "/"), "/api/traces")
	var err error
	traceProvider, err = tracing.TracerProvider(1, jaegerEndpoint, "mw-api-svc", "dev")
	ctx := context.Background()
	if err != nil {
		logger.Log.Warn(ctx, "unable to start tracing")
	}
}
