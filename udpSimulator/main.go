package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	tracesdk "go.opentelemetry.io/otel/sdk/trace"
	"proj.com/udp/messageHandlers"
	"proj.com/udp/tracing"
)

const (
	serverAddr   = "127.0.0.1"
	messageCount = 10
	servicePort  = ":3000"
)

var traceProvider *tracesdk.TracerProvider

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
	jaegerEndpoint := fmt.Sprintf("%s%s", strings.TrimRight("http://simplest-collector.default.svc.cluster.local:14268", "/"), "/api/traces")
	var err error
	ctx := context.Background()
	traceProvider, err = tracing.TracerProvider(ctx, 1, jaegerEndpoint, "udp-svc", "dev")
	if err != nil {
		log.Println("unable to start tracing")
	}
}
