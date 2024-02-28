package messageHandlers

import (
	"fmt"
	"log"

	"github.com/IBM/sarama"
)

// CreateTopic ex
func CreateTopic(topicName string) {
	// Configure the Kafka producer
	config := sarama.NewConfig()
	config.Net.SASL.Enable = true
	config.Net.SASL.User = "user1"
	config.Net.SASL.Password = "G8wOkt9OZc"
	config.Producer.Return.Successes = true

	// Create the Kafka broker addresses
	brokerAddresses := []string{
		"kafka.default.svc.cluster.local:9092",
	}

	// Create a new Kafka producer
	producer, err := sarama.NewSyncProducer(brokerAddresses, config)
	if err != nil {
		log.Fatal("Error creating Kafka producer:", err)
	}
	defer func() {
		if err := producer.Close(); err != nil {
			log.Fatal("Error closing Kafka producer:", err)
		}
	}()

	// Define the topic configuration
	topicDetail := &sarama.TopicDetail{
		NumPartitions:     1,
		ReplicationFactor: 1,
	}

	// Create topics on Kafka
	admin, err := sarama.NewClusterAdmin(brokerAddresses, config)
	if err != nil {
		log.Fatal("Error creating Kafka admin:", err)
	}
	defer func() {
		if err := admin.Close(); err != nil {
			log.Fatal("Error closing Kafka admin:", err)
		}
	}()

	if exists, err := topicExists(brokerAddresses, config, topicName); err != nil {
		log.Fatalf("Error checking topic existence: %v", err)
	} else if !exists {
		err = admin.CreateTopic(topicName, topicDetail, false)
		if err != nil {
			log.Fatal("Error creating Kafka topic:", err)
		}
	}

	fmt.Println("Kafka topic created successfully!")
}

func topicExists(brokerAddresses []string, config *sarama.Config, topic string) (bool, error) {
	admin, err := sarama.NewClusterAdmin(brokerAddresses, config)
	if err != nil {
		return false, err
	}
	defer func() {
		if err := admin.Close(); err != nil {
			log.Printf("Error closing Kafka admin: %v", err)
		}
	}()

	topics, err := admin.ListTopics()
	if err != nil {
		return false, err
	}

	_, exists := topics[topic]
	return exists, nil
}
