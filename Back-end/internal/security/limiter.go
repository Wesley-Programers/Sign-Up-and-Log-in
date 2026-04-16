package security

import (
	"context"
	"fmt"
	"time"
	
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

	script := `
		local current = redis.call("INCR", KEYS[1])
		if current == 1 then
			redis.call("PEXPIRE", KEYS[1], ARGV[1])
		end
		return current`

	result, err := rl.client.Eval(ctx, script, []string{fullKey}, window.Milliseconds()).Result()
	if err != nil {
		return false, err
	}

	count := result.(int64)
	return int(count) <= maxAttempts, nil
}

func (rl *RedisLimiter) ResetLimit(ctx context.Context, key string) error {
	fullKey := fmt.Sprintf("rl:%s", key)
	return rl.client.Del(ctx, fullKey).Err()
}