package main

import (
	"context"
	"flag"
	"io"
	"log"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	lpb "github.com/exsmund/grpc-psql-example/proto/logger"
	upb "github.com/exsmund/grpc-psql-example/proto/user"
)

var serverAddr = flag.String("serveraddr", "localhost:8001", "the address to connect to")
var loggerAddr = flag.String("loggeraddr", "localhost:8002", "the address to connect to")

func saveUser(c upb.UserRepoClient, ctx context.Context, u *upb.User) *upb.User {
	log.Printf("Save User: {%s, %s}\n", u.GetName(), u.GetEmail())
	saveRes, err := c.SaveUser(ctx, &upb.SaveUserRequest{Data: u})
	if err != nil {
		log.Fatalf("Could not save: %v", err)
	}
	log.Printf("Status: %s, ID: %d\n", saveRes.GetStatus(), saveRes.GetData().GetId())
	return saveRes.GetData()
}

func deleteUser(c upb.UserRepoClient, ctx context.Context, ID int32) {
	delRes, err := c.DeleteUser(ctx, &upb.DeleteUserRequest{Id: ID})
	log.Printf("Delete user with ID: %d", ID)
	if err != nil {
		log.Fatalf("Could not delete: %v", err)
	}
	log.Printf("Status: %s", delRes.GetStatus())
}

func getUsers(c upb.UserRepoClient, ctx context.Context) {
	log.Printf("Getting users")
	stream, err := c.GetUsers(ctx, &upb.GetUsersRequest{})
	if err != nil {
		log.Fatalf("Could not get users: %v", err)
	}
	for {
		u, err := stream.Recv()
		if err == io.EOF {
			break
		}
		if err != nil {
			log.Fatalf("Could not get users: %v", err)
		}
		log.Printf("User: %d, %s, %s", u.GetId(), u.GetName(), u.GetEmail())
	}
}

func getLogs(c lpb.LoggerRepoClient, ctx context.Context) {
	log.Printf("Getting log")
	logs, err := c.GetTail(ctx, &lpb.GetTailRequest{})
	if err != nil {
		log.Fatalf("Could not get log: %v", err)
	}
	for _, row := range logs.Logs {
		log.Printf("%s: %s", row.GetTs(), row.GetMsg())
	}
}

func main() {
	flag.Parse()
	time.Sleep(time.Second * 10)
	// Set up a connection to the server.
	sConn, err := grpc.Dial(*serverAddr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	defer sConn.Close()
	c := upb.NewUserRepoClient(sConn)
	log.Printf("Connected gRPC %s", *serverAddr)
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
	defer cancel()

	// Test saving used data
	u := &upb.User{Name: "Ivan", Email: "example@email1"}
	u = saveUser(c, ctx, u)

	time.Sleep(time.Second)
	log.Print("--------")
	// If last request was more than a minute ago,
	// we will recieve all users include last saved.
	getUsers(c, ctx)

	time.Sleep(time.Second)
	log.Print("--------")
	// Test saving used data.
	u2 := &upb.User{Name: "Vasya", Email: "example@email2"}
	u2 = saveUser(c, ctx, u2)

	time.Sleep(time.Second)
	log.Print("--------")
	// Last user was created less then a minute ago,
	// so we will not recieve it
	getUsers(c, ctx)

	time.Sleep(time.Second)
	log.Print("--------")
	// Test deleting used by gotten IDs
	deleteUser(c, ctx, u.Id)
	deleteUser(c, ctx, u2.Id)

	// Test deleting used by unknown ID
	deleteUser(c, ctx, 0)

	time.Sleep(time.Second)
	log.Print("--------")
	// Set up a connection to the logger.
	lConn, err := grpc.Dial(*loggerAddr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	defer lConn.Close()
	l := lpb.NewLoggerRepoClient(lConn)
	log.Printf("Connected gRPC %s", *loggerAddr)
	// If last request was more than a minute ago,
	// we will recieve all users include last saved.
	ctx2, cancel2 := context.WithTimeout(context.Background(), time.Second)
	defer cancel2()
	getLogs(l, ctx2)
}
