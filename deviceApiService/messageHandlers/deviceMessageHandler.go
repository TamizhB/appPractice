package messageHandlers

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"strings"
)

// Send test
func Send(topic string, payload string) (int, string, error) {
	//kafkaService := "kafka-service.default.svc.cluster.local:9099"
	client := &http.Client{}
	req, err := http.NewRequest("POST", fmt.Sprintf("http://kafka-service.default.svc.cluster.local:9099/kafkasvc/msg/send"), strings.NewReader(payload))
	if err != nil {
		fmt.Println("Error creating HTTP request:", err)
		return http.StatusInternalServerError, "", err
	}

	req.Header.Set("message-topic", topic)
	// Make the HTTP request
	resp, err := client.Do(req)
	if err != nil {
		log.Println("Error making HTTP request:", err)
		return http.StatusInternalServerError, "", err
	}
	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Println("Error reading response body:", err)
		return http.StatusInternalServerError, "", err
	}
	return resp.StatusCode, string(body), nil
}

// // Receive test
// func Receive(ctx *gin.Context) {
// 	//kafkaService := "kafka-service.default.svc.cluster.local:9099"
// 	client := &http.Client{}
// 	req, err := http.NewRequest("GET", fmt.Sprintf("http://kafka-service.default.svc.cluster.local:9099/kafkasvc/msg/receive"), nil)
// 	if err != nil {
// 		fmt.Println("Error creating HTTP request:", err)
// 		return
// 	}

// 	// Make the HTTP request
// 	resp, err := client.Do(req)
// 	if err != nil {
// 		log.Println("Error making HTTP request:", err)
// 		return
// 	}
// 	defer resp.Body.Close()
// 	body, err := io.ReadAll(resp.Body)
// 	if err != nil {
// 		log.Println("Error reading response body:", err)
// 		return
// 	}

// 	fmt.Println("Response:", string(body))

// }
