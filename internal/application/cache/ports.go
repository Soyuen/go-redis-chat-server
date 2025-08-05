// internal/application/cache/ports.go
//go:generate mockgen -destination=mocks/mock_redis_cache.go -package=mocks github.com/Soyuen/go-redis-chat-server/internal/application/cache RedisCache

package cache

import "context"

type RedisCache interface {
	Set(ctx context.Context, key string, value interface{}) error
	Get(ctx context.Context, key string) (string, error)
	Delete(ctx context.Context, key string) error

	// ZSet operations for chatroom user management
	ZAdd(ctx context.Context, key, member string, score float64) error
	ZRem(ctx context.Context, key, member string) error
	ZCard(ctx context.Context, key string) (int64, error)
	ZRange(ctx context.Context, key string, start, stop int64) ([]string, error)
}
