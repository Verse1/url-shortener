package db

import (
	"github.com/go-redis/redis/v8"
	"context"
	"os"
)

var Ctxt=context.Background()

func DBinit(db int) *redis.Client {
	rdb := redis.NewClient(&redis.Options{
		Addr: os.Getenv("REDIS_HOST"),
		Password: os.Getenv("REDIS_PASSWORD"),
		DB: db,
	})
	return rdb
}

