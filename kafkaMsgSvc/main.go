package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	tracesdk "go.opentelemetry.io/otel/sdk/trace"
	"proj.com/kafkasvc/messageHandlers"
	"proj.com/kafkasvc/tracing"
)

const (
	servicePort = ":9099"
)

var traceProvider *tracesdk.TracerProvider

func main() {
	router := gin.New()
	base := router.Group("/kafkasvc")

	server := &http.Server{Addr: servicePort, Handler: router}

	// create default topic
	messageHandlers.CreateTopic("configRequest")
	messageHandlers.CreateTopic("configResponse")

	go messageHandlers.ConsumeAll()

	msg := base.Group("msg")
	msg.POST("/send", messageHandlers.Send)
	//msg.GET("/receive", messageHandlers.Receive)

	if err := server.ListenAndServe(); err != nil {
		panic(err)
	}
}

func initTracing() {
	// init tracing
	jaegerEndpoint := fmt.Sprintf("%s%s", strings.TrimRight("http://simplest-collector.default.svc.cluster.local:14268", "/"), "/api/traces")
	var err error
	ctx := context.Background()
	traceProvider, err = tracing.TracerProvider(ctx, 1, jaegerEndpoint, "kafka-svc", "dev")
	if err != nil {
		log.Println("unable to start tracing")
	}
}
