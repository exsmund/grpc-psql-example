package main

import (
	"bytes"
	"context"
	"database/sql"
	"encoding/gob"
	"fmt"
	"log"
	"time"

	pb "github.com/exsmund/grpc-psql-example/proto/user"
	"github.com/go-redis/redis"
)

type userServer struct {
	pb.UnimplementedUserRepoServer
	db    *sql.DB
	redis *redis.Client
	ctx   context.Context
}

type StoredUser struct {
	Id    int32
	Name  string
	Email string
}

type StoredUsers struct {
	Users []StoredUser
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
	redisUsers, ok := s.getUsersFromRedis()
	if ok {
		log.Print("Retrieve users from Redis")
		for _, u := range *redisUsers {
			stream.Send(&pb.User{Id: u.Id, Name: u.Name, Email: u.Email})
		}
		return nil
	}

	log.Print("Retrieve users from DB")
	rows, err := s.db.Query("SELECT id, email, name FROM users")
	if err != nil {
		return err
	}
	defer rows.Close()

	// Loop through rows, using Scan to assign column data to struct fields.
	var users []StoredUser
	for rows.Next() {
		u := &pb.User{}
		if err := rows.Scan(&u.Id, &u.Email, &u.Name); err != nil {
			return fmt.Errorf("Failed while scaning row: %s", err)
		}
		users = append(users, StoredUser{Id: u.Id, Email: u.Email, Name: u.Name})
		stream.Send(u)
	}
	if err := rows.Err(); err != nil {
		return fmt.Errorf("Failed getting users: %s", err)
	}

	log.Print("Store to Redis")
	s.storeUsersIntoRedis(&users)

	return nil
}

func (s *userServer) getUsersFromRedis() (*[]StoredUser, bool) {
	fromRedis, err := s.redis.Get(redisUsersKey).Bytes()
	if err == nil {
		buf := bytes.NewBuffer([]byte(fromRedis))
		dec := gob.NewDecoder(buf)
		var s StoredUsers
		err = dec.Decode(&s)
		if err == nil {
			return &s.Users, true
		} else {
			log.Printf("Error: getting users from Redis: %s", err)
		}
	}
	return nil, false
}

func (s *userServer) storeUsersIntoRedis(users *[]StoredUser) {
	forRedis := &StoredUsers{Users: *users}
	var buf bytes.Buffer
	enc := gob.NewEncoder(&buf)
	err := enc.Encode(forRedis)
	if err == nil {
		err = s.redis.Set(redisUsersKey, buf.Bytes(), time.Minute).Err()
		if err != nil {
			log.Printf("Error: storing users into Redis: %s", err)
		}
	} else {
		log.Printf("Error: storing users into Redis: %s", err)
	}
}

func newServer(db *sql.DB, redis *redis.Client) *userServer {
	s := &userServer{
		db:    db,
		redis: redis,
		ctx:   context.Background(),
	}
	return s
}
