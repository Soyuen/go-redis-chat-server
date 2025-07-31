package redis

import (
	"context"
	"fmt"

	"github.com/Soyuen/go-redis-chat-server/internal/config"

	"github.com/redis/go-redis/v9"
)

type RedisAdapter struct {
	client *redis.Client
}

func NewRedisAdapter(cfg config.RedisConfig) (*RedisAdapter, error) {
	client := redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%s:%d", cfg.Host(), cfg.Port()),
		Password: cfg.Password(),
		DB:       cfg.DB(),
	})
	return &RedisAdapter{client: client}, nil
}

func (r *RedisAdapter) RawClient() *redis.Client {
	return r.client
}

func (r *RedisAdapter) Set(ctx context.Context, key string, value interface{}) error {
	return r.client.Set(ctx, key, value, 0).Err()
}

func (r *RedisAdapter) Get(ctx context.Context, key string) (string, error) {
	return r.client.Get(ctx, key).Result()
}

func (r *RedisAdapter) Delete(ctx context.Context, key string) error {
	return r.client.Del(ctx, key).Err()
}

// ZAdd adds a member with a score to a sorted set.
func (r *RedisAdapter) ZAdd(ctx context.Context, key, member string, score float64) error {
	z := &redis.Z{Score: score, Member: member}
	return r.client.ZAdd(ctx, key, *z).Err()
}

// ZRem removes a member from a sorted set.
func (r *RedisAdapter) ZRem(ctx context.Context, key, member string) error {
	return r.client.ZRem(ctx, key, member).Err()
}

// ZCard returns the cardinality (number of elements) of the sorted set.
func (r *RedisAdapter) ZCard(ctx context.Context, key string) (int64, error) {
	return r.client.ZCard(ctx, key).Result()
}

// ZRange returns the members in the sorted set within the given range.
func (r *RedisAdapter) ZRange(ctx context.Context, key string, start, stop int64) ([]string, error) {
	return r.client.ZRange(ctx, key, start, stop).Result()
}
