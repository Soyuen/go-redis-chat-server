package chat

//go:generate mockgen -destination=mocks/mock_chat_interface.go -package=mocks github.com/Soyuen/go-redis-chat-server/internal/application/chat ChatService

import (
	"context"

	domainchat "github.com/Soyuen/go-redis-chat-server/internal/domain/chat"
)

type ChatService interface {
	ProcessIncoming(raw []byte, sender, channel string) (*domainchat.Message, error)
	CreateRoom(roomName string) error
	BroadcastSystemMessage(ctx context.Context, channel, nickname, action string) error
	JoinChannel(ctx context.Context, channel, nickname string) error
}
