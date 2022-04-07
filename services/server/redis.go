package main

import "github.com/go-redis/redis"

func connectRedis(addr string, pass string, db int) *redis.Client {
	rdb := redis.NewClient(&redis.Options{
		Addr:     addr,
		Password: pass,
		DB:       db,
	})
	return rdb
}
