package chat

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"

	"github.com/Soyuen/go-redis-chat-server/internal/application/realtime"
	domainchat "github.com/Soyuen/go-redis-chat-server/internal/domain/chat"
	"github.com/Soyuen/go-redis-chat-server/internal/presenter"
)

type chatService struct {
	channelManager realtime.ChannelManager
	redisSub       realtime.ChannelEventSubscriber
	presenter      presenter.MessagePresenterInterface
	goFunc         func(func())
	memberRepo     domainchat.ChatMemberRepository
}

func NewChatService(channelManager realtime.ChannelManager, redisSub realtime.ChannelEventSubscriber, presenter presenter.MessagePresenterInterface, memberRepo domainchat.ChatMemberRepository) ChatService {
	return &chatService{
		channelManager: channelManager,
		redisSub:       redisSub,
		presenter:      presenter,
		goFunc:         func(f func()) { go f() },
		memberRepo:     memberRepo,
	}
}

func (s *chatService) CreateRoom(roomName string) error {
	s.channelManager.GetOrCreateChannel(roomName)
	s.goFunc(func() {
		s.redisSub.Start(roomName)
	})
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

func (s *chatService) AddUserToRoom(ctx context.Context, room, user string) error {
	return s.memberRepo.AddUserToRoom(ctx, room, user)
}
