package quiz

import (
	"context"
	"fmt"
	"github.com/redis/go-redis/v9"
)

type redisAdapter struct {
	client *redis.Client
}

func (re *redisAdapter) BindCode(quiz Quiz) error {
	return re.client.Set(context.Background(), fmt.Sprintf("quiz-%s", quiz.Code), quiz.Id, 0).Err()
}
func (re *redisAdapter) UnbindCode(code string) error {
	return re.client.Del(context.Background(), fmt.Sprintf("quiz-%s", code)).Err()
}
func (re *redisAdapter) GetQuiz(code string) (string, error) {
	return re.client.Get(context.Background(), fmt.Sprintf("quiz-%s", code)).Result()
}
