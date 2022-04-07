package main

import (
	"context"
	"database/sql"
	"fmt"
	"log"

	pb "github.com/exsmund/grpc-psql-example/proto/user"
)

type userServer struct {
	pb.UnimplementedUserRepoServer
	db *sql.DB
}

func (s *userServer) SaveUser(ctx context.Context, req *pb.SaveUserRequest) (*pb.SaveUserResponse, error) {
	user := req.GetData()
	email := user.GetEmail()
	name := user.GetName()
	log.Printf("Save User: %v, %v", name, email)
	var id int32
	sqlStatement := `
		INSERT INTO users (email, name)
		VALUES ($1, $2)
		RETURNING id
	`
	err := s.db.QueryRow(sqlStatement, email, name).Scan(&id)
	if err != nil {
		return nil, fmt.Errorf("Failed user saving: %v", err)
	}
	log.Printf("User Id: %v", id)
	user.Id = id
	return &pb.SaveUserResponse{Status: pb.Status_Success, Data: user}, nil
}

func (s *userServer) DeleteUser(ctx context.Context, req *pb.DeleteUserRequest) (*pb.DeleteUserResponse, error) {
	id := req.GetId()
	if id <= 0 {
		return &pb.DeleteUserResponse{Status: pb.Status_RequestError}, nil
	}
	log.Printf("Delete User by id: %v", id)
	sqlStatement := `
		DELETE FROM users
		WHERE id = $1;`
	res, err := s.db.Exec(sqlStatement, id)
	if err != nil {
		return nil, err
	}
	c, err := res.RowsAffected()
	if err != nil {
		return nil, err
	}
	if c > 0 {
		return &pb.DeleteUserResponse{Status: pb.Status_Success}, nil
	}
	return &pb.DeleteUserResponse{Status: pb.Status_NothinToDo}, nil
}

func (s *userServer) GetUsers(req *pb.GetUsersRequest, stream pb.UserRepo_GetUsersServer) error {
	rows, err := s.db.Query("SELECT id, email, name FROM users")
	if err != nil {
		return err
	}
	defer rows.Close()

	// Loop through rows, using Scan to assign column data to struct fields.
	for rows.Next() {
		u := &pb.User{}
		if err := rows.Scan(&u.Id, &u.Email, &u.Name); err != nil {
			return fmt.Errorf("Failed while scaning row: %s", err)
		}
		stream.Send(u)
	}
	if err := rows.Err(); err != nil {
		return fmt.Errorf("Failed getting users: %s", err)
	}
	return nil
}

func newServer(db *sql.DB) *userServer {
	s := &userServer{db: db}
	return s
}
