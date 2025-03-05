package quizzes

import (
	"context"
	"fmt"
	"github.com/redis/go-redis/v9"
)

type redisAdapter struct {
	client *redis.Client
}

func (re *redisAdapter) BindCode(ownerId string, quiz Quiz) error {
	return re.client.Set(context.Background(), quiz.Code, fmt.Sprintf("%s@%s", ownerId, quiz.Id), 0).Err()
}
func (re *redisAdapter) UnbindCode(code string) error {
	return re.client.Del(context.Background(), code).Err()
}
func (re *redisAdapter) GetQuiz(code string) (string, error) {
	return re.client.Get(context.Background(), code).Result()
}
