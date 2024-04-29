package session

import (
	"chat-bot/pkg/database/redis"
	"context"
	"log"
	"time"

	redisdriver "github.com/go-redis/redis/v8"
	"github.com/google/uuid"
)

// GenerateSessionID creates a new unique session ID
func GenerateSessionID() string {
	return uuid.New().String()
}

// Set sets a key-value pair in a specific session
func Set(ctx context.Context, sessionID, key, value string) {
	fullKey := "session:" + sessionID + ":" + key
	err := redis.RDB.Set(ctx, fullKey, value, time.Minute*5).Err()
	if err != nil {
		log.Print("Error Setting information on Redis: ", err.Error())
	}
}

// Get retrieves a value by key from a specific session
func Get(ctx context.Context, sessionID, key string) string {
	fullKey := "session:" + sessionID + ":" + key
	val, err := redis.RDB.Get(ctx, fullKey).Result()
	if err == redisdriver.Nil {
		log.Print("[DEBUB] Key does not exist: ", fullKey)
		return ""
	} else if err != nil {
		log.Print("Error Getting information from Redis: ", err.Error())
		return ""
	}
	return val
}
