package messageHandlers

import (
	"context"
	"log"

	kafka "github.com/segmentio/kafka-go"
)

func main() {
	// Define the Kafka broker address
	broker := []string{"localhost:9092"}

	// Create a new Kafka reader
	reader := kafka.NewReader(kafka.ReaderConfig{
		Brokers:   broker,
		Topic:     "my-topic",
		Partition: 0,
		MinBytes:  10e3, // 10KB
		MaxBytes:  10e6, // 10MB
	})

	// Create a context for canceling the consumer
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// Handle received messages
	for {
		message, err := reader.FetchMessage(ctx)
		if err != nil {
			log.Fatalf("Error reading message: %v", err)
		}

		log.Printf("Received message: %s\n", string(message.Value))

		// Mark the message as processed
		reader.CommitMessages(ctx, message)
	}
}
