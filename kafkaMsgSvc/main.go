package main

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"proj.com/kafkasvc/messageHandlers"
)

const (
	servicePort = ":9099"
)

func main() {

	// //ctx := context.Background()
	// brokerAddresses := []string{"kafka.default.svc.cluster.local:9092"}

	// config := sarama.NewConfig()
	// config.Net.SASL.Enable = true
	// config.Net.SASL.User = "user1"
	// config.Net.SASL.Password = "G8wOkt9OZc" // Replace with the actual password
	// config.Producer.Return.Successes = true

	router := gin.New()
	base := router.Group("/api/v1/device")

	server := &http.Server{Addr: servicePort, Handler: router}

	messageHandlers.CreateTopic()
	msg := base.Group("msgTest")
	msg.POST("/send", messageHandlers.Send)
	msg.GET("/receive", messageHandlers.Receive)

	if err := server.ListenAndServe(); err != nil {
		panic(err)
	}
}
