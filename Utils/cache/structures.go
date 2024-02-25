package cache

import (
	"github.com/redis/go-redis/v9"
)

var RClient RedisClient

type RedisClient struct {
	client *redis.Client
}
