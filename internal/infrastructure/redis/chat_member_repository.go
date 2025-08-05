package redis

import (
	"context"
	"errors"
	"time"

	"github.com/Soyuen/go-redis-chat-server/internal/application/cache"
	"github.com/Soyuen/go-redis-chat-server/internal/domain/chat"
	"github.com/go-redis/redis"
)

type ChatMemberRepository struct {
	redis cache.RedisCache
}

func NewChatMemberRepository(redis cache.RedisCache) *ChatMemberRepository {
	return &ChatMemberRepository{redis: redis}
}

func (r *ChatMemberRepository) AddUserToRoom(ctx context.Context, room, user string) error {
	// Use current Unix timestamp as score
	return r.redis.ZAdd(ctx, room, user, float64(time.Now().Unix()))
}

func (r *ChatMemberRepository) RemoveUserFromRoom(ctx context.Context, room, user string) error {
	return r.redis.ZRem(ctx, room, user)
}

func (r *ChatMemberRepository) GetRoomUserCount(ctx context.Context, room string) (int64, error) {
	return r.redis.ZCard(ctx, room)
}

func (r *ChatMemberRepository) GetRoomUserList(ctx context.Context, room string) ([]string, error) {
	// Get all members in the room
	return r.redis.ZRange(ctx, room, 0, -1)
}

func (r *ChatMemberRepository) UserExists(ctx context.Context, room, user string) (bool, error) {
	_, err := r.redis.ZScore(ctx, room, user)
	if err != nil {
		if errors.Is(err, redis.Nil) {
			return false, nil
		}
		return false, err
	}
	return true, nil
}

// Ensure type implements the domain interface
var _ chat.ChatMemberRepository = (*ChatMemberRepository)(nil)
