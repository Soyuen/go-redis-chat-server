package chat

import (
	"encoding/json"
	"errors"

	domainchat "github.com/Soyuen/go-redis-chat-server/internal/domain/chat"

	"github.com/Soyuen/go-redis-chat-server/pkg/realtimeiface"
)

type ChatService interface {
	CreateRoom(roomName string) error
	ProcessIncoming(raw []byte, sender, channel string) (*domainchat.Message, error)
}

type chatService struct {
	channelManager realtimeiface.ChatChannelManager
	redisSub       realtimeiface.ChannelEventSubscriber
}

func NewChatService(channelManager realtimeiface.ChatChannelManager, redisSub realtimeiface.ChannelEventSubscriber) ChatService {
	return &chatService{
		channelManager: channelManager,
		redisSub:       redisSub,
	}
}

func (s *chatService) CreateRoom(roomName string) error {
	s.channelManager.GetOrCreateChannel(roomName)
	go s.redisSub.Start(roomName)
	return nil
}

func (s *chatService) ProcessIncoming(raw []byte, sender, channel string) (*domainchat.Message, error) {
	var payload struct {
		Content string `json:"message"`
	}
	if err := json.Unmarshal(raw, &payload); err != nil {
		return nil, errors.New("invalid message format")
	}

	return domainchat.NewMessage(sender, channel, payload.Content)
}
