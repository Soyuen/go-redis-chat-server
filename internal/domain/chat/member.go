package chat

import "context"

type ChatMemberRepository interface {
	AddUserToRoom(ctx context.Context, room, user string) error
	RemoveUserFromRoom(ctx context.Context, room, user string) error
	GetRoomUserCount(ctx context.Context, room string) (int64, error)
	GetRoomUserList(ctx context.Context, room string) ([]string, error)
}
