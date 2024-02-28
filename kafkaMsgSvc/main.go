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
