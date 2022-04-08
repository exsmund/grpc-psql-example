package main

import (
	"log"

	"github.com/confluentinc/confluent-kafka-go/kafka"
)

type Consumer struct {
	c  *kafka.Consumer
	ch *CH
}

func (c *Consumer) Destroy() {
	c.c.Close()
}

func (c *Consumer) Subscribe(topic string) {
	c.c.SubscribeTopics([]string{topic}, nil)
	for {
		msg, err := c.c.ReadMessage(-1)
		if err == nil {
			log.Printf("Message on %s: %s\n", msg.TopicPartition, string(msg.Value))
			c.ch.Write(string(msg.Value))
		} else {
			// The client will automatically try to recover from all errors.
			log.Printf("Consumer error: %v (%v)\n", err, msg)
		}
	}
}

func NewConsumer(server, kafkagroup string, ch *CH) (*Consumer, error) {
	c, err := kafka.NewConsumer(&kafka.ConfigMap{
		"bootstrap.servers": server,
		"group.id":          kafkagroup,
		"auto.offset.reset": "earliest",
	})
	if err != nil {
		return nil, err
	}

	cons := &Consumer{
		c,
		ch,
	}
	return cons, nil
}
