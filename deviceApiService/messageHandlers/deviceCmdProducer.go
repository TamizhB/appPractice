package messageHandlers

import (
	"context"
	"log"

	kafka "github.com/segmentio/kafka-go"
)

func main() {
	// Define the Kafka broker address
	broker := []string{"localhost:9092"}

	// Create a new Kafka writer
	writer := kafka.NewWriter(kafka.WriterConfig{
		Brokers:  broker,
		Topic:    "my-topic",
		Balancer: &kafka.LeastBytes{},
	})

	// Message to send
	message := kafka.Message{
		Key:   []byte("key"),
		Value: []byte("Hello, Kafka!"),
	}

	// Send the message
	err := writer.WriteMessages(context.Background(), message)
	if err != nil {
		log.Fatalf("Failed to send message: %v", err)
	}

	log.Println("Message sent successfully")

	// Close the writer
	writer.Close()
}
