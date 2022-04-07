package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"net"

	"google.golang.org/grpc"

	pb "github.com/exsmund/grpc-psql-example/proto"
)

var (
	port = flag.Int("port", 8001, "The server port")
)

type userServer struct {
	pb.UnimplementedUserRepoServer
}

func (s *userServer) SaveUser(ctx context.Context, user *pb.User) (*pb.SaveUserResponse, error) {
	log.Printf("Save User: %v, %v", user.Email, user.Name)
	user.Id = []byte("a")
	return &pb.SaveUserResponse{Status: pb.Status_Success, Data: user}, nil
}

func (s *userServer) DeleteUser(ctx context.Context, uuid *pb.UserUUID) (*pb.DeleteUserResponse, error) {
	if string(uuid.Id) == "a" {
		log.Printf("Delete User with UUID %v", uuid)
		return &pb.DeleteUserResponse{Status: pb.Status_Success}, nil
	}
	return &pb.DeleteUserResponse{Status: pb.Status_NotFound}, nil
}

func (s *userServer) GetUsers(req *pb.GetUsersRequest, stream pb.UserRepo_GetUsersServer) error {
	users := []*pb.User{
		{Name: "Ivan", Email: ""},
		{Name: "Vasya", Email: ""},
	}
	for _, u := range users {
		stream.Send(u)
	}
	return nil
}

func newServer() *userServer {
	s := &userServer{}
	return s
}

func main() {
	flag.Parse()
	lis, err := net.Listen("tcp", fmt.Sprintf("localhost:%d", *port))
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	grpcServer := grpc.NewServer()
	pb.RegisterUserRepoServer(grpcServer, newServer())
	grpcServer.Serve(lis)
}
