package redis

import (
	"context"
	"fmt"

	"github.com/go-redis/redis/v8"
)

var RDB *redis.Client

// Initialize connects to Redis and should be called at the beginning of the application.
func Initialize() {
	RDB = redis.NewClient(&redis.Options{
		Addr:     "redis-server:6379", // Docker service name
		Password: "",                  // No password set
		DB:       0,                   // Use default DB
	})

	ctx := context.Background()
	// Check the connection
	pong, err := RDB.Ping(ctx).Result()
	if err != nil {
		fmt.Println("Redis connect error:", err)
	} else {
		fmt.Println("Redis connected:", pong)
	}
}
