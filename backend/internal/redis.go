package internal

import (
	"context"
	"log"

	"github.com/go-redis/redis/v8"
	rstore "github.com/rbcervilla/redisstore/v8"
)

var store *rstore.RedisStore

func LoadRedis() {
	var err error

	client := redis.NewClient(&redis.Options{
		Addr: "localhost:6379",
	})

	store, err = rstore.NewRedisStore(context.Background(), client)
	if err != nil {
		log.Fatal("failed to create redis store: ", err)
	}

	store.KeyPrefix("session_token")
}
