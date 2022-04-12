package main

import (
	"context"
	"log"
	"time"

	"github.com/confluentinc/confluent-kafka-go/kafka"
)

type KafkaConsumer struct {
	ok    bool
	c     *kafka.Consumer
	topic string
	ch    *CH
}

func (c *KafkaConsumer) Destroy() {
	c.c.Close()
}

func (k *KafkaConsumer) Connect(server, kafkagroup string, ch *CH, topic string) {
	var err error
	defer func() {
		k.ok = err == nil
	}()

	k.topic = topic
	k.ch = ch

	c, err := kafka.NewConsumer(&kafka.ConfigMap{
		"bootstrap.servers":  server,
		"group.id":           kafkagroup,
		"enable.auto.commit": "false",
	})
	if err != nil {
		return
	}

	k.c = c

	err = k.createTopic(topic)
	if err != nil {
		return
	}

	err = k.c.SubscribeTopics([]string{topic}, nil)
	if err != nil {
		return
	}

}

func (k *KafkaConsumer) createTopic(topic string) error {
	a, err := kafka.NewAdminClientFromConsumer(k.c)
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

func (k *KafkaConsumer) Listen() {
	for {
		msg, err := k.c.ReadMessage(-1)
		if err == nil {
			k.handleMsg(msg)
		} else {
			// The client will automatically try to recover from all errors.
			log.Printf("Consumer error: %v (%v)\n", err, msg)
		}
	}
}

func (k *KafkaConsumer) handleMsg(msg *kafka.Message) {
	log.Printf("Message on %s: %s\n", msg.TopicPartition, string(msg.Value))
	err := k.ch.Write(string(msg.Value))
	if err == nil {
		log.Printf("Commit the Message")
		k.c.CommitMessage(msg)
	}
}

func (k *KafkaConsumer) Close() {
	k.c.Close()
}
