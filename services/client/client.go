package main

import (
	"context"
	"flag"
	"io"
	"log"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	pb "github.com/exsmund/grpc-psql-example/proto/user"
)

var addr = flag.String("addr", "localhost:8001", "the address to connect to")

func saveUser(c pb.UserRepoClient, ctx context.Context, u *pb.User) *pb.User {
	log.Printf("Save User: {%s, %s}\n", u.GetName(), u.GetEmail())
	saveRes, err := c.SaveUser(ctx, &pb.SaveUserRequest{Data: u})
	if err != nil {
		log.Fatalf("Could not save: %v", err)
	}
	log.Printf("Status: %s, ID: %d\n", saveRes.GetStatus(), saveRes.GetData().GetId())
	return saveRes.GetData()
}

func deleteUser(c pb.UserRepoClient, ctx context.Context, ID int32) {
	delRes, err := c.DeleteUser(ctx, &pb.DeleteUserRequest{Id: ID})
	log.Printf("Delete user with ID: %d", ID)
	if err != nil {
		log.Fatalf("Could not delete: %v", err)
	}
	log.Printf("Status: %s", delRes.GetStatus())
}

func getUsers(c pb.UserRepoClient, ctx context.Context) {
	log.Printf("Getting users")
	stream, err := c.GetUsers(ctx, &pb.GetUsersRequest{})
	if err != nil {
		log.Fatalf("Could not gt users: %v", err)
	}
	for {
		u, err := stream.Recv()
		if err == io.EOF {
			break
		}
		if err != nil {
			log.Fatalf("Could not gt users: %v", err)
		}
		log.Printf("User: %d, %s, %s", u.GetId(), u.GetName(), u.GetEmail())
	}
}

func main() {
	flag.Parse()
	// Set up a connection to the server.
	conn, err := grpc.Dial(*addr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	defer conn.Close()
	c := pb.NewUserRepoClient(conn)

	// Contact the server and print out its response.
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	// Test saving used data
	u := &pb.User{Name: "Ivan", Email: "example@email1"}
	u = saveUser(c, ctx, u)
	// Test saving used data
	u2 := &pb.User{Name: "Vasya", Email: "example@email2"}
	u2 = saveUser(c, ctx, u2)

	getUsers(c, ctx)

	// Test deleting used by gotten IDs
	deleteUser(c, ctx, u.Id)
	deleteUser(c, ctx, u2.Id)

	// Test deleting used by unknown ID
	deleteUser(c, ctx, 0)

}
