package chat

import (
	domainchat "github.com/Soyuen/go-redis-chat-server/internal/domain/chat"
)

type ChatService interface {
	ProcessIncoming(raw []byte, sender, channel string) (*domainchat.Message, error)
	CreateRoom(roomName string) error
	BroadcastSystemMessage(channel, nickname, action string) error
}
