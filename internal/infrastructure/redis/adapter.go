package redis

import (
	"context"
	"fmt"

	pkgRedis "github.com/Soyuen/go-redis-chat-server/pkg/redis"

	"github.com/redis/go-redis/v9"
)

type RedisAdapter struct {
	client *redis.Client
}

func NewRedisAdapter(cfg pkgRedis.Config) (*RedisAdapter, error) {
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
