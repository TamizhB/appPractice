package messageHandlers

import (
	"log"
	"os"
	"os/signal"
	"time"

	"github.com/IBM/sarama"
)

// ConsumeAll
func ConsumeAll() {
	log.Printf("Kafka consumer started to receive messages...")

	brokerAddresses := []string{"kafka.default.svc.cluster.local:9092"}

	config := sarama.NewConfig()
	config.Net.SASL.Enable = true
	config.Net.SASL.User = "user1"
	config.Net.SASL.Password = "G8wOkt9OZc" // Replace with the actual password
	config.Producer.Return.Successes = true

	consumer, err := sarama.NewConsumer(brokerAddresses, config)
	if err != nil {
		log.Fatalf("Error creating Kafka consumer: %v", err)
	}
	defer func() {
		if err := consumer.Close(); err != nil {
			log.Fatalf("Error closing Kafka consumer: %v", err)
		}
	}()

	topics, err := consumer.Topics()
	if err != nil {
		log.Fatalf("Error getting list of topics: %v", err)
	}

	// Subscribe to all topics
	for _, topic := range topics {
		partitionConsumer, err := consumer.ConsumePartition(topic, 0, sarama.OffsetNewest)
		if err != nil {
			log.Printf("Error subscribing to Kafka topic %s: %v", topic, err)
			continue
		}
		defer func() {
			if err := partitionConsumer.Close(); err != nil {
				log.Printf("Error closing Kafka partition consumer for topic %s: %v", topic, err)
			}
		}()

		// Handle incoming messages in a separate goroutine
		go func(topic string) {
			for {
				select {
				case msg := <-partitionConsumer.Messages():
					HandleMessage(msg)
				case err := <-partitionConsumer.Errors():
					log.Printf("Error on topic %s: %v\n", topic, err)
				}
			}
		}(topic)
	}

	// Wait for a signal to gracefully shut down the consumer
	sigchan := make(chan os.Signal, 1)
	signal.Notify(sigchan, os.Interrupt)
	<-sigchan

	log.Println("Closing Kafka consumer...")
}

func HandleMessage(msg *sarama.ConsumerMessage) {
	if "configRequest" == msg.Topic {
		log.Printf("Received config request %s at %v.. Action:To be logged\n", string(msg.Value), time.Now())
	}
	if "configResponse" == msg.Topic {
		log.Printf("Received config response from device %s at %v.. Action:To be logged\n", string(msg.Value), time.Now())
	}
}
