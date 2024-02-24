package main

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"proj.com/udp/messageHandlers"
)

const (
	port         = 8888
	serverAddr   = "127.0.0.1"
	messageCount = 10
	servicePort  = ":3000"
)

func main() {
	go messageHandlers.Serve()

	router := gin.New()
	base := router.Group("/udpSvc")

	server := &http.Server{Addr: servicePort, Handler: router}

	messageHandlers.CreateTopic()
	msg := base.Group("msgTest")
	msg.POST("/send", messageHandlers.sendMessage)

	if err := server.ListenAndServe(); err != nil {
		panic(err)
	}

	time.Sleep(5 * time.Second)

}

func Serve() {

}
