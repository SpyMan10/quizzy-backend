package services

import (
	"github.com/redis/go-redis/v9"
	"quizzy.app/backend/quizzy/cfg"
)

func ConfigureRedis(cfg cfg.AppConfig) *redis.Client {
	opt, err := redis.ParseURL(cfg.RedisUri)
	if err != nil {
		panic(err)
	}

	return redis.NewClient(opt)
}
