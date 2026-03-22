package security

import (
	"context"
	"fmt"
	"time"
	
	// "net/http"
	"github.com/redis/go-redis/v9"
)

type RedisLimiter struct {
	client *redis.Client
}

func NewRedisLimiter(addr string) *RedisLimiter {
	return &RedisLimiter{
		client: redis.NewClient(&redis.Options{
			Addr: addr,
		}),
	}
}

func (rl *RedisLimiter) CheckLimit(ctx context.Context, key string, maxAttempts int, window time.Duration) (bool, error) {
	fullKey := fmt.Sprintf("rl:%s", key)

	count, err := rl.client.Incr(ctx, fullKey).Result()
	if err != nil {
		return false, err
	}

	if count == 1 {
		rl.client.Expire(ctx, fullKey, window)
	}

	return int(count) <= maxAttempts, nil
}