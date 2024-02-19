package engine

import (
	"fmt"
	"log"
	"os"

	"github.com/confluentinc/confluent-kafka-go/kafka"
)

func CreateConsumer() *kafka.Consumer {
	fmt.Println("Creating consumer")
	config := &kafka.ConfigMap{
		"bootstrap.servers": os.Getenv("KAFKA_BOOTSTRAP_SERVERS"),
		"group.id":          os.Getenv("KAFKA_GROUP_ID"),
		"auto.offset.reset": "earliest",
	}
	c, err := kafka.NewConsumer(config)
	if err != nil {
		log.Fatal(err)
	}
	return c
}

func CreateProducer() *kafka.Producer {
	fmt.Println("Creating producer")
	config := &kafka.ConfigMap{
		"bootstrap.servers": os.Getenv("KAFKA_BOOTSTRAP_SERVERS"),
	}
	p, err := kafka.NewProducer(config)
	if err != nil {
		log.Fatal(err)
	}
	return p
}
