package main

import (
	"flag"
	"fmt"
	"log"
	"net"

	"google.golang.org/grpc"

	pb "github.com/exsmund/grpc-psql-example/proto/user"
)

var (
	port    = flag.Int("port", 8001, "The server port")
	dbport  = flag.Int("dbport", 5432, "The PostgreSQL port")
	dbhost  = flag.String("dbaddr", "127.0.0.1", "The PostgreSQL host")
	dbname  = flag.String("dbname", "example", "The PostgreSQL database name")
	dbuname = flag.String("dbuname", "postgres", "The PostgreSQL username")
	dbpass  = flag.String("dbpass", "example", "The PostgreSQL password")
)

func main() {
	flag.Parse()

	var conninfo string = fmt.Sprintf("host=%s port=%d dbname=%s user=%s password=%s sslmode=disable", *dbhost, *dbport, *dbname, *dbuname, *dbpass)
	db, err := connectDB(conninfo)
	if err != nil {
		log.Fatalf("Failed DB connection: %v", err)
	}
	defer db.Close()
	fmt.Println("Connected DB")

	lis, err := net.Listen("tcp", fmt.Sprintf("localhost:%d", *port))
	if err != nil {
		log.Fatalf("Failed to listen: %v", err)
	}
	grpcServer := grpc.NewServer()
	pb.RegisterUserRepoServer(grpcServer, newServer(db))
	fmt.Println("Start server")
	grpcServer.Serve(lis)
}
