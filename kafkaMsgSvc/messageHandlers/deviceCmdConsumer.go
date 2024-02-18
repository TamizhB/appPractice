package messageHandlers

import (
	"log"

	"github.com/IBM/sarama"
	"github.com/gin-gonic/gin"
)

// Receive test
func Receive(ctx *gin.Context) {
	log.Printf("Receiving message")

	// Define the Kafka broker addresses
	brokerAddresses := []string{"kafka.default.svc.cluster.local:9092"}

	// Configure the Kafka consumer
	config := sarama.NewConfig()
	config.Net.SASL.Enable = true
	config.Net.SASL.User = "user1"
	config.Net.SASL.Password = "G8wOkt9OZc" // Replace with the actual password
	config.Producer.Return.Successes = true

	// Create a new Kafka consumer
	consumer, err := sarama.NewConsumer(brokerAddresses, config)
	if err != nil {
		log.Fatalf("Error creating Kafka consumer: %v", err)
	}
	defer func() {
		if err := consumer.Close(); err != nil {
			log.Fatalf("Error closing Kafka consumer: %v", err)
		}
	}()

	// Define the Kafka topic
	topic := "my-topic"

	// Subscribe to the topic
	partitionConsumer, err := consumer.ConsumePartition(topic, 0, sarama.OffsetNewest)
	if err != nil {
		log.Fatalf("Error subscribing to Kafka topic: %v", err)
	}
	defer func() {
		if err := partitionConsumer.Close(); err != nil {
			log.Fatalf("Error closing Kafka partition consumer: %v", err)
		}
	}()

	// Handle received messages
	for {
		select {
		case message := <-partitionConsumer.Messages():
			log.Printf("Received message: %s\n", string(message.Value))
			//defer consumer.Close()
		case err := <-partitionConsumer.Errors():
			log.Printf("Error receiving message: %v", err)
		case <-ctx.Done():
			log.Println("Consumer stopped")
			return
		}
	}
}
