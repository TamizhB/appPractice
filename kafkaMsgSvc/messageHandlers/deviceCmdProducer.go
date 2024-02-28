package messageHandlers

import (
	"context"
	"fmt"
	"io"
	"log"
	"net/http"
	"strings"

	"github.com/IBM/sarama"
	"github.com/gin-gonic/gin"
	"go.opentelemetry.io/otel"
)

// Send test
func Send(ctx *gin.Context) {
	tr := otel.Tracer("udp-service")
	ctx1 := context.Background()
	ctx1, span := tr.Start(ctx1, "udp-svc")
	defer span.End()
	log.Printf("Sending message")

	reqBody, err := io.ReadAll(ctx.Request.Body)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Error reading request body"})
		return
	}

	brokerAddresses := []string{"kafka.default.svc.cluster.local:9092"}

	// Configure the Kafka producer
	config := sarama.NewConfig()
	config.Net.SASL.Enable = true
	config.Net.SASL.User = "user1"
	config.Net.SASL.Password = "G8wOkt9OZc"
	config.Producer.Return.Successes = true
	// Create a new Kafka producer
	producer, err := sarama.NewSyncProducer(brokerAddresses, config)
	if err != nil {
		log.Fatalf("Failed to create Kafka producer: %v", err)
	}
	defer func() {
		if err := producer.Close(); err != nil {
			log.Fatalf("Error closing Kafka producer: %v", err)
		}
	}()

	// Define the Kafka topic
	topic := ctx.Request.Header.Get("message-topic")

	// Create the topic if it does not exist
	admin, err := sarama.NewClusterAdmin(brokerAddresses, config)
	if err != nil {
		log.Fatalf("Error creating Kafka admin: %v", err)
	}
	defer func() {
		if err := admin.Close(); err != nil {
			log.Fatalf("Error closing Kafka admin: %v", err)
		}
	}()

	// Message to send
	message := &sarama.ProducerMessage{
		Topic: topic,
		Key:   sarama.StringEncoder("config"),
		Value: sarama.StringEncoder(string(reqBody)),
	}

	// Send the message
	_, _, err = producer.SendMessage(message)
	if err != nil {
		log.Fatalf("Failed to send message: %v", err)
	}

	log.Println("Message pushed to kafka successfully")

	log.Println("Sending udp request")
	respStatus, respBody, err := SendConfigRequest(string(reqBody))

	if err != nil {
		ctx.JSON(respStatus, err.Error())
		return
	}
	ctx.JSON(respStatus, respBody)
}

// SendRequest
func SendConfigRequest(payload string) (int, string, error) {
	//kafkaService := "kafka-service.default.svc.cluster.local:9099"
	client := &http.Client{}
	req, err := http.NewRequest("POST", fmt.Sprintf("http://udp-service.default.svc.cluster.local:3000/udpSvc/msg/send"), strings.NewReader(payload))
	if err != nil {
		fmt.Println("Error creating HTTP request:", err)
		return http.StatusInternalServerError, "", err
	}

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
	log.Println("Response from udp server:", string(body))
	return resp.StatusCode, string(body), nil
}
