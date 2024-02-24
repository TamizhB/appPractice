package messageHandlers

import (
	"fmt"
	"io"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
)

// Send test
func Send(ctx *gin.Context) {
	//kafkaService := "kafka-service.default.svc.cluster.local:9099"
	client := &http.Client{}
	req, err := http.NewRequest("POST", fmt.Sprintf("http://kafka-service.default.svc.cluster.local:9099/api/v1/device/msgTest/send"), nil)
	if err != nil {
		fmt.Println("Error creating HTTP request:", err)
		return
	}

	// Make the HTTP request
	resp, err := client.Do(req)
	if err != nil {
		log.Println("Error making HTTP request:", err)
		return
	}
	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Println("Error reading response body:", err)
		return
	}

	fmt.Println("Response:", string(body))

}

// Receive test
func Receive(ctx *gin.Context) {
	//kafkaService := "kafka-service.default.svc.cluster.local:9099"
	client := &http.Client{}
	req, err := http.NewRequest("GET", fmt.Sprintf("http://kafka-service.default.svc.cluster.local:9099/api/v1/device/msgTest/receive"), nil)
	if err != nil {
		fmt.Println("Error creating HTTP request:", err)
		return
	}

	// Make the HTTP request
	resp, err := client.Do(req)
	if err != nil {
		log.Println("Error making HTTP request:", err)
		return
	}
	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Println("Error reading response body:", err)
		return
	}

	fmt.Println("Response:", string(body))

}
