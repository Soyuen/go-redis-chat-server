package chat

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"

	"github.com/Soyuen/go-redis-chat-server/internal/application/realtime"
	domainchat "github.com/Soyuen/go-redis-chat-server/internal/domain/chat"
	apperr "github.com/Soyuen/go-redis-chat-server/internal/errors"
	"github.com/Soyuen/go-redis-chat-server/internal/presenter"
	"github.com/Soyuen/go-redis-chat-server/pkg/loggeriface"
)

type chatService struct {
	channelManager realtime.ChannelManager
	redisSub       realtime.ChannelEventSubscriber
	presenter      presenter.MessagePresenterInterface
	memberRepo     domainchat.ChatMemberRepository
	logger         loggeriface.Logger
	goFunc         func(func())
}

func NewChatService(channelManager realtime.ChannelManager, redisSub realtime.ChannelEventSubscriber, presenter presenter.MessagePresenterInterface, memberRepo domainchat.ChatMemberRepository, logger loggeriface.Logger) ChatService {
	return &chatService{
		channelManager: channelManager,
		redisSub:       redisSub,
		presenter:      presenter,
		memberRepo:     memberRepo,
		logger:         logger,
		goFunc:         func(f func()) { go f() },
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

func (s *chatService) BroadcastSystemMessage(ctx context.Context, channel, nickname, action string) error {
	content := fmt.Sprintf("%s %s the chat.", nickname, action)
	count, err := s.memberRepo.GetRoomUserCount(ctx, channel)
	if err != nil {
		return err
	}
	payload := map[string]interface{}{
		"content": content,
		"count":   count,
	}
	payloadBytes, _ := json.Marshal(payload)
	return s.createAndBroadcastMessage("System", channel, string(payloadBytes))
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

// 原子性 JoinChannel 實作
func (s *chatService) JoinChannel(ctx context.Context, channel, nickname string) error {
	// 1. Create room
	if err := s.CreateRoom(channel); err != nil {
		return fmt.Errorf("%w: %v", apperr.ErrCreateRoom, err)
	}
	// 2. Add user
	if err := s.memberRepo.AddUserToRoom(ctx, channel, nickname); err != nil {
		return fmt.Errorf("%w: %v", apperr.ErrAddUserToRoom, err)
	}
	// 3. Broadcast
	if err := s.BroadcastSystemMessage(ctx, channel, nickname, "joined"); err != nil {
		if rmErr := s.memberRepo.RemoveUserFromRoom(ctx, channel, nickname); rmErr != nil {
			s.logger.Warnw("compensate RemoveUserFromRoom failed", "err", rmErr, "channel", channel, "nickname", nickname)
		}
		return fmt.Errorf("%w: %v", apperr.ErrBroadcastSystemMessage, err)
	}
	return nil
}

func (s *chatService) UserExists(ctx context.Context, room, user string) (bool, error) {
	return s.memberRepo.UserExists(ctx, room, user)
}
