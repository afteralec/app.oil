package configs

import (
	"os"

	redis "github.com/redis/go-redis/v9"
)

func Redis() redis.Options {
	return redis.Options{
		Addr:     os.Getenv("REDIS_ADDR"),
		Password: "",
		DB:       0,
		Protocol: 3,
	}
}
