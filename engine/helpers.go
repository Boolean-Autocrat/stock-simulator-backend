package engine

import (
	"log"
	"os"

	"github.com/confluentinc/confluent-kafka-go/kafka"
)

func CreateConsumer() *kafka.Consumer {
	config := &kafka.ConfigMap{
		"bootstrap.servers": os.Getenv("KAFKA_BOOTSTRAP_SERVERS"),
		"group.id":          os.Getenv("KAFKA_GROUP_ID"),
		"auto.offset.reset": "earliest",
	}
	c, err := kafka.NewConsumer(config)
	if err != nil {
		log.Fatal(err)
	}
	err = c.Subscribe("orders", nil)
	if err != nil {
		log.Fatal(err)
	}
	log.Println("Consumer created, subscribed to orders")
	return c
}

func CreateProducer() *kafka.Producer {
	config := &kafka.ConfigMap{
		"bootstrap.servers": os.Getenv("KAFKA_BOOTSTRAP_SERVERS"),
	}
	p, err := kafka.NewProducer(config)
	if err != nil {
		log.Fatal(err)
	}
	return p
}
