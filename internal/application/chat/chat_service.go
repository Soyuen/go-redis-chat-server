package chat

import (
	"github.com/Soyuen/go-redis-chat-server/pkg/realtimeiface"
)

type ChatService interface {
	CreateRoom(roomName string) error
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
