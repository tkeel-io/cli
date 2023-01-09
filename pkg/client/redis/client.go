package redis

import (
	"fmt"

	"github.com/go-redis/redis/v8"
)

func NewClient(host string, port int, password string, db int) *redis.Client {
	rdb := redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%s:%d", host, port),
		Password: password, // no password set
		DB:       db,       // use default DB
	})
	return rdb
}
