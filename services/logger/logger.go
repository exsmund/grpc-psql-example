package main

import (
	"flag"
	"fmt"
	"log"
	"net"
	"sync"

	"google.golang.org/grpc"

	"github.com/exsmund/grpc-psql-example/packages/color"
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

	var wg sync.WaitGroup

	ch := &CH{}
	wg.Add(1)
	go func() {
		defer wg.Done()
		ch.Connect(*CHAddr, *CHDB, *CHUser, *CHPass)
	}()

	kafka := &KafkaConsumer{}
	wg.Add(1)
	go func() {
		defer wg.Done()
		kafka.Connect(*kafkaServer, *kafkaGroup, ch, kafkaTopic)
	}()

	wg.Wait()

	if ch.ok {
		log.Println("Connection to the Clickhouse: " + color.Green + "done" + color.Reset)
		defer ch.Close()
	} else {
		log.Fatalf(color.Red+"ERROR"+color.Reset+": Failed Clickhouse connection (%s:%d)", CHAddr)
	}

	if kafka.ok {
		log.Println("Connection to the Kafka: " + color.Green + "done" + color.Reset)
		defer kafka.Close()
		log.Print("Starting consuming Kafka: " + color.Green + "done" + color.Reset)
		go kafka.Listen()
	} else {
		log.Fatalf(color.Red+"ERROR"+color.Reset+": Failed Kafka connection (%s)", *kafkaServer)
	}

	lis, err := net.Listen("tcp", fmt.Sprintf("0.0.0.0:%d", *port))
	if err != nil {
		log.Fatalf("Failed to listen: %v", err)
	}
	defer lis.Close()
	grpcServer := grpc.NewServer()
	lpb.RegisterLoggerRepoServer(grpcServer, newServer(ch))
	log.Printf("Start gRPC server %s "+color.Green+"done"+color.Reset, lis.Addr().String())
	defer grpcServer.Stop()
	grpcServer.Serve(lis)

	// kafka, err := NewConsumer(*kafkaServer, *kafkaGroup, ch)
	// if err != nil {
	// 	log.Fatal(err)
	// }
	// defer kafka.Destroy()
	// log.Print("Start consuming Kafka: OK")
	// kafka.Subscribe(kafkaTopic)
}
