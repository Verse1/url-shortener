package db

import (
	"github.com/go-redis/redis/v8"
	"context"
)

var ctxt=context.Background()

func DBinit(db int) *redis.Client {
	rdb := redis.NewClient(&redis.Options{
		Addr: "localhost:6379",
		DB: db,
	})
	return rdb
}

