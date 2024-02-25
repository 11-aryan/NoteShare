package cache

import (
	"context"
	"time"

	"github.com/gofiber/fiber/v2/log"
	"github.com/redis/go-redis/v9"
)

func InitRedis() {
	client := redis.NewClient(
		&redis.Options{
			Addr: "localhost:6379",
			DB:   0,
		},
	)
	ctx, cancel := context.WithTimeout(context.TODO(), 10*time.Second)
	defer cancel()
	err := client.Ping(ctx).Err()
	if err != nil {
		log.Error(err)
		return
	}
	RClient = RedisClient{client}
}
