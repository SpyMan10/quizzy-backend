package services

import (
	"github.com/redis/go-redis/v9"
	"quizzy.app/backend/quizzy/cfg"
)

func ConfigureRedis(cfg cfg.AppConfig) (*redis.Client, error) {
	if opt, err := redis.ParseURL(cfg.RedisUri); err != nil {
		return nil, err
	} else {
		return redis.NewClient(opt), nil
	}
}
