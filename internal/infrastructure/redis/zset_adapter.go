package redis

import (
	"context"

	appredis "github.com/Soyuen/go-redis-chat-server/internal/application/redis"
	"github.com/redis/go-redis/v9"
)

type ZSetAdapter struct {
	client *redis.Client
}

var _ appredis.ZSetRepository = (*ZSetAdapter)(nil)

func NewZSetAdapter(client *redis.Client) *ZSetAdapter {
	return &ZSetAdapter{client: client}
}

func (r *ZSetAdapter) ZAdd(ctx context.Context, key string, score float64, member string) error {
	z := &redis.Z{
		Score:  score,
		Member: member,
	}
	return r.client.ZAdd(ctx, key, *z).Err()
}

func (r *ZSetAdapter) ZRem(ctx context.Context, key string, member string) error {
	return r.client.ZRem(ctx, key, member).Err()
}

func (r *ZSetAdapter) ZRange(ctx context.Context, key string, start, stop int64) ([]string, error) {
	return r.client.ZRange(ctx, key, start, stop).Result()
}

func (r *ZSetAdapter) ZScore(ctx context.Context, key string, member string) (float64, error) {
	return r.client.ZScore(ctx, key, member).Result()
}
