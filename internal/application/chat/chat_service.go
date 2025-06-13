package chat

import (
	"encoding/json"
	"errors"
	"fmt"

	domainchat "github.com/Soyuen/go-redis-chat-server/internal/domain/chat"
	"github.com/Soyuen/go-redis-chat-server/internal/presenter"

	"github.com/Soyuen/go-redis-chat-server/pkg/realtimeiface"
)

type ChatService interface {
	CreateRoom(roomName string) error
	ProcessIncoming(raw []byte, sender, channel string) (*domainchat.Message, error)
	BroadcastSystemMessage(channel, nickname, action string) error
}

type chatService struct {
	channelManager realtimeiface.ChatChannelManager
	redisSub       realtimeiface.ChannelEventSubscriber
	presenter      presenter.MessagePresenter
}

func NewChatService(channelManager realtimeiface.ChatChannelManager, redisSub realtimeiface.ChannelEventSubscriber, presenter presenter.MessagePresenter) ChatService {
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

func (s *chatService) BroadcastSystemMessage(channel, nickname, action string) error {
	content := fmt.Sprintf("%s %s the chat.", nickname, action)
	return s.createAndBroadcastMessage("System", channel, content)
}

func (s *chatService) createAndBroadcastMessage(sender, channel, content string) error {
	msg, err := domainchat.NewMessage(sender, channel, content)
	if err != nil {
		return err
	}

	if formatted := s.presenter.Format(msg); formatted != nil {
		s.channelManager.Broadcast(*formatted)
	}

	return nil
}
