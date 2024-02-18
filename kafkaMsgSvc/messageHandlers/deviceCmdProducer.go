package messageHandlers

import (
	"log"

	"github.com/IBM/sarama"
	"github.com/gin-gonic/gin"
)

// Send test
func Send(ctx *gin.Context) {
	log.Printf("Sending message")

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
	topic := "my-topic"

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
		Key:   sarama.StringEncoder("key"),
		Value: sarama.StringEncoder("Hello, Kafka!"),
	}

	// Send the message
	_, _, err = producer.SendMessage(message)
	if err != nil {
		log.Fatalf("Failed to send message: %v", err)
	}

	log.Println("Message sent successfully")
}
