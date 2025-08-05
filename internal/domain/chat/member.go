package chat

//go:generate mockgen -destination=mocks/mock_chat_member_repository.go -package=mocks github.com/Soyuen/go-redis-chat-server/internal/domain/chat ChatMemberRepository

import "context"

type ChatMemberRepository interface {
	AddUserToRoom(ctx context.Context, room, user string) error
	RemoveUserFromRoom(ctx context.Context, room, user string) error
	GetRoomUserCount(ctx context.Context, room string) (int64, error)
	GetRoomUserList(ctx context.Context, room string) ([]string, error)
	UserExists(ctx context.Context, room, user string) (bool, error)
}
