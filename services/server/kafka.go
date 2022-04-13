package main

import (
	"context"
	"log"
	"time"

	"github.com/confluentinc/confluent-kafka-go/kafka"
)

type KafkaProducer struct {
	ok    bool
	p     *kafka.Producer
	topic string
	flush int
}

func (k *KafkaProducer) Connect(conf *KafkaConfig) {
	// log.Printf("create kafka produser with server %s and topic %s", conf.server, conf.topic)
	var err error
	defer func() {
		k.ok = err == nil
	}()

	k.topic = conf.topic
	k.flush = 1000

	p, err := kafka.NewProducer(&kafka.ConfigMap{
		"bootstrap.servers": conf.server,
		"acks":              0, //  The producer will not wait for any acknowledgment from the server at all
	})
	if err != nil {
		return
	}
	k.p = p

	err = k.createTopic(conf.topic)
	// go m.deliveryHandler()
}

func (k *KafkaProducer) createTopic(topic string) error {
	a, err := kafka.NewAdminClientFromProducer(k.p)
	if err != nil {
		return err
	}
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	results, err := a.CreateTopics(
		ctx,
		[]kafka.TopicSpecification{{
			Topic:             topic,
			NumPartitions:     1,
			ReplicationFactor: 1}},
		kafka.SetAdminOperationTimeout(time.Second*5))
	if err != nil {
		log.Printf("Admin Client request error: %v\n", err)
		return err
	}
	for _, result := range results {
		if result.Error.Code() != kafka.ErrNoError && result.Error.Code() != kafka.ErrTopicAlreadyExists {
			log.Printf("Failed to create topic: %v\n", result.Error)
			return result.Error
		}
		log.Printf("Creation topic: %v\n", result)
	}
	return nil
}

func (k *KafkaProducer) Close() {
	k.p.Close()
}

// Delivery report handler for produced messages
func (k *KafkaProducer) deliveryHandler() {
	for e := range k.p.Events() {
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

func (p *KafkaProducer) Produce(msg string) {
	log.Printf("Producing message: %v", msg)
	p.p.Produce(&kafka.Message{
		TopicPartition: kafka.TopicPartition{Topic: &p.topic, Partition: kafka.PartitionAny},
		Value:          []byte(msg),
	}, nil)
	p.p.Flush(p.flush)
}
