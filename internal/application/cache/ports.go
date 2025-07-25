// internal/application/cache/ports.go
package cache

import "context"

type RedisCache interface {
	Set(ctx context.Context, key string, value interface{}) error
	Get(ctx context.Context, key string) (string, error)
	Delete(ctx context.Context, key string) error
}
