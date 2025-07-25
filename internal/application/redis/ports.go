package redis

import "context"

type ZSetRepository interface {
	ZAdd(ctx context.Context, key string, score float64, member string) error
	ZRem(ctx context.Context, key string, member string) error
	ZRange(ctx context.Context, key string, start, stop int64) ([]string, error)
	ZScore(ctx context.Context, key string, member string) (float64, error)
}
