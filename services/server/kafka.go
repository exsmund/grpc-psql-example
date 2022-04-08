package main

import (
	"log"

	"github.com/confluentinc/confluent-kafka-go/kafka"
)

type Producer struct {
	p     *kafka.Producer
	topic string
	flush int
}

func (p *Producer) Destroy() {
	p.p.Close()
}

// Delivery report handler for produced messages
func (p *Producer) deliveryHandler() {
	for e := range p.p.Events() {
		switch ev := e.(type) {
		case *kafka.Message:
			if ev.TopicPartition.Error != nil {
				log.Printf("Delivery failed: %v\n", ev.TopicPartition)
			} else {
				log.Printf("Delivered message to %v\n", ev.TopicPartition)
			}
		}
	}
}

func (p *Producer) Produce(msg string) {
	p.p.Produce(&kafka.Message{
		TopicPartition: kafka.TopicPartition{Topic: &p.topic, Partition: kafka.PartitionAny},
		Value:          []byte(msg),
	}, nil)
	p.p.Flush(p.flush)
}

func NewProducer(server string, topic string) (*Producer, error) {
	log.Printf("create kafka produser with server %s and topic %s", server, topic)
	p, err := kafka.NewProducer(&kafka.ConfigMap{"bootstrap.servers": server})
	if err != nil {
		return nil, err
	}

	prod := &Producer{
		p,
		topic,
		1000,
	}

	go prod.deliveryHandler()

	return prod, nil
}
