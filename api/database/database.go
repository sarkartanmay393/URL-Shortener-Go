package database

import (
	"context"
	"fmt"
	"github.com/go-redis/redis/v8"
	"log"
	"os"
)

var Ctx = context.Background()

func CreateClient(dbNo int) *redis.Client {
	rdbURL := fmt.Sprintf("redis://%s:%s@localhost:6379/%d", "redis", os.Getenv("DB_PASS"), dbNo)
	opt, err := redis.ParseURL(rdbURL)
	if err != nil {
		log.Printf("Error creating redis client: %v\n", err)
	}

	return redis.NewClient(opt)
}
