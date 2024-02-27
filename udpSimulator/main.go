package main

import (
	"net/http"

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

	messageHandlers.SendMessage(nil)

	server := &http.Server{Addr: servicePort, Handler: router}
	msg := base.Group("msgTest")
	msg.POST("/send", messageHandlers.SendMessage)

	if err := server.ListenAndServe(); err != nil {
		panic(err)
	}
}
