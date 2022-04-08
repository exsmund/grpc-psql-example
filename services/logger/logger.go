package main

import (
	"flag"
	"fmt"
	"log"
	"net"

	"google.golang.org/grpc"

	lpb "github.com/exsmund/grpc-psql-example/proto/logger"
)

var (
	port        = flag.Int("port", 8002, "The server port")
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
	rows, err := ch.Read()
	if err != nil {
		log.Fatal(err)
	}
	for _, r := range *rows {
		log.Print(r)
	}

	lis, err := net.Listen("tcp", fmt.Sprintf("localhost:%d", *port))
	if err != nil {
		log.Fatalf("Failed to listen: %v", err)
	}
	defer lis.Close()
	grpcServer := grpc.NewServer()
	lpb.RegisterLoggerRepoServer(grpcServer, newServer(ch))
	log.Printf("Start gRPC server %s", lis.Addr().String())
	defer grpcServer.Stop()
	go grpcServer.Serve(lis)

	kafka, err := NewConsumer(*kafkaServer, *kafkaGroup, ch)
	if err != nil {
		log.Fatal(err)
	}
	defer kafka.Destroy()
	log.Print("Start consuming Kafka: OK")
	kafka.Subscribe(kafkaTopic)
}
