package quizzes

import (
	"context"
	"fmt"
	"github.com/redis/go-redis/v9"
)

type RedisCodeResolver struct {
	client *redis.Client
}

func (re *RedisCodeResolver) BindCode(ownerId string, quiz Quiz) error {
	return re.client.Set(context.Background(), quiz.Code, fmt.Sprintf("%store@%store", ownerId, quiz.Id), 0).Err()
}

func (re *RedisCodeResolver) UnbindCode(code string) error {
	return re.client.Del(context.Background(), code).Err()
}

func (re *RedisCodeResolver) GetQuiz(code string) (string, error) {
	return re.client.Get(context.Background(), code).Result()
}
