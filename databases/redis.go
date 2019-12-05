package databases

import "github.com/go-redis/redis"

var RedisSession *redis.Client

func init() {
	RedisSession = redis.NewClient(&redis.Options{
		Addr:     "redis:6379",
		Password: "",
		DB:       0,
	})
}