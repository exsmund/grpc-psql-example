package main

import (
	"flag"
	"log"

	"github.com/confluentinc/confluent-kafka-go/kafka"
)

var (
	kafkaserver = flag.String("kafkaserver", "localhost:9092", "The Kafka server")
	kafkagroup  = flag.String("kafkagroup", "myGroup", "The Kafka group id")
)

func main() {
	flag.Parse()

	c, err := kafka.NewConsumer(&kafka.ConfigMap{
		"bootstrap.servers": *kafkaserver,
		"group.id":          *kafkagroup,
		"auto.offset.reset": "earliest",
	})
	log.Print("Start consume Kafka")
	defer c.Close()

	if err != nil {
		panic(err)
	}

	c.SubscribeTopics([]string{"log", "^aRegex.*[Tt]opic"}, nil)

	for {
		msg, err := c.ReadMessage(-1)
		if err == nil {
			log.Printf("Message on %s: %s\n", msg.TopicPartition, string(msg.Value))
		} else {
			// The client will automatically try to recover from all errors.
			log.Printf("Consumer error: %v (%v)\n", err, msg)
		}
	}

}
