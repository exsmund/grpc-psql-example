package main

import (
	"flag"
	"log"
)

var (
	kafkaServer = flag.String("kafkaserver", "localhost:9092", "The Kafka server")
	kafkaGroup  = flag.String("kafkagroup", "myGroup", "The Kafka group id")
	CHAddr      = flag.String("chaddr", "127.0.0.1:9000", "The Clickhouse address")
	CHDB        = flag.String("chdb", "log", "The Clickhouse DB")
	CHUser      = flag.String("chuser", "default", "The Clickhouse user")
	CHPass      = flag.String("chpass", "", "The Clickhouse password")
)

func main() {
	flag.Parse()

	ch, err := newCH(*CHAddr, *CHDB, *CHUser, *CHPass)
	if err != nil {
		log.Fatal(err)
	}
	log.Print("Connect Clickhouse: OK")

	kafka, err := NewConsumer(*kafkaServer, *kafkaGroup, ch)
	if err != nil {
		log.Fatal(err)
	}
	defer kafka.Destroy()
	log.Print("Start consuming Kafka: OK")
	kafka.Subscribe(kafkaTopic)
}
