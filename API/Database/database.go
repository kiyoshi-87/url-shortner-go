package database

import (
	"os"

	"github.com/redis/go-redis/v9"
)

func CreateClient(dbNo int) *redis.Client {
	client := redis.NewClient(&redis.Options{
		Addr:     os.Getenv("DB_ADDRS"),
		Password: os.Getenv("DB_PASSWD"),
		DB:       dbNo,
	})

	return client
}
