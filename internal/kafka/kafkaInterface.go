package kafka

import (
	"bookorder/internal/datatypes"
	"bookorder/internal/db"
	"encoding/json"
	"fmt"
	"github.com/IBM/sarama"
	"log"
	"time"
)

func InitKafka() {
	config := sarama.NewConfig()
	consumer, err := sarama.NewConsumer([]string{"127.0.0.1:9092"}, config)
	if err != nil {
		log.Fatal("Error creating Kafka consumer:", err)
	}
	defer func() {
		if err := consumer.Close(); err != nil {
			log.Fatal("Error closing Kafka consumer:", err)
		}
	}()

	topic := "Register"   // Replace with the actual Kafka topic name
	partition := int32(0) // Replace with the partition you want to consume from

	// Create a new consumer for the specified topic and partition
	consumerPartition, err := consumer.ConsumePartition(topic, partition, sarama.OffsetOldest)
	if err != nil {
		log.Fatal("Error creating Kafka partition consumer:", err)
	}
	defer func() {
		if err := consumerPartition.Close(); err != nil {
			log.Fatal("Error closing Kafka partition consumer:", err)
		}
	}()

	fmt.Printf("Listening to messages on topic '%s'...\n", topic)

	// Continuously listen for and print Kafka messages
	for {
		select {
		case msg := <-consumerPartition.Messages():
			fmt.Printf("Received message: %s\n", string(msg.Value))
			var order datatypes.Order
			if err := json.Unmarshal(msg.Value, &order); err != nil {
				log.Println("Error parsing Kafka message:", err)
			} else {
				order.CreatedAt = time.Now()
				savingErr := db.SaveOrder(&order)
				if savingErr != nil {
					fmt.Println("Can't save data to database. Error message : ", savingErr)
				} else {
					fmt.Println("Success: Data saved to database successfully")
				}

			}
		case err := <-consumerPartition.Errors():
			log.Println("Kafka consumer error:", err)
		}
	}
}
