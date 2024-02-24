package messageHandlers

import (
	"log"
	"os"
	"os/signal"

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

	// Handle incoming messages in a separate goroutine
	go func() {
		for {
			select {
			case msg := <-partitionConsumer.Messages():
				log.Printf("Received message from topic %s: %s\n", msg.Topic, string(msg.Value))
			case err := <-partitionConsumer.Errors():
				log.Printf("Error: %v\n", err)
			}
		}
	}()

	// Wait for a signal to gracefully shut down the consumer
	sigchan := make(chan os.Signal, 1)
	signal.Notify(sigchan, os.Interrupt)
	<-sigchan

	log.Println("Closing Kafka consumer...")
}
