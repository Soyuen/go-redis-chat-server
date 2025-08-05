package chat

import (
	"context"
	"encoding/json"
	"testing"

	realtimemock "github.com/Soyuen/go-redis-chat-server/internal/application/realtime/mocks"
	domainchatmock "github.com/Soyuen/go-redis-chat-server/internal/domain/chat/mocks"
	logmock "github.com/Soyuen/go-redis-chat-server/internal/infrastructure/logger/mocks"
	presentermock "github.com/Soyuen/go-redis-chat-server/internal/presenter/mocks"

	"github.com/Soyuen/go-redis-chat-server/internal/testhelper"

	"github.com/Soyuen/go-redis-chat-server/internal/application/realtime"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
)

func TestChatService_CreateRoom(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockCM := realtimemock.NewMockChannelManager(ctrl)
	mockSub := realtimemock.NewMockChannelEventSubscriber(ctrl)
	mockPresenter := presentermock.NewMockMessagePresenterInterface(ctrl)
	mockMemberRepo := domainchatmock.NewMockChatMemberRepository(ctrl)
	mockLogger := logmock.NewMockLogger(ctrl)

	mockCM.EXPECT().GetOrCreateChannel(testhelper.ChannelTest).Times(1)
	mockSub.EXPECT().Start(testhelper.ChannelTest).Times(1)

	service := NewChatService(mockCM, mockSub, mockPresenter, mockMemberRepo, mockLogger)
	service.(*chatService).goFunc = func(f func()) { f() }

	err := service.CreateRoom(testhelper.ChannelTest)
	assert.NoError(t, err)
}

func TestChatService_ProcessIncoming_Valid(t *testing.T) {
	service := NewChatService(nil, nil, presentermock.NewMockMessagePresenterInterface(nil), nil, nil)

	raw, _ := json.Marshal(map[string]string{testhelper.KeyMessage: testhelper.MessageHello})
	msg, err := service.ProcessIncoming(raw, testhelper.SenderAlice, testhelper.ChannelTest)

	assert.NoError(t, err)
	assert.Equal(t, testhelper.SenderAlice, msg.Sender)
	assert.Equal(t, testhelper.ChannelTest, msg.Channel)
	assert.Equal(t, testhelper.MessageHello, msg.Content)
}

func TestChatService_ProcessIncoming_InvalidJSON(t *testing.T) {
	service := NewChatService(nil, nil, presentermock.NewMockMessagePresenterInterface(nil), nil, nil)

	_, err := service.ProcessIncoming([]byte("invalid json"), testhelper.SenderAlice, testhelper.ChannelTest)

	assert.EqualError(t, err, "invalid message format")
}

func TestChatService_JoinChannel(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	ctx := context.Background()
	room := testhelper.ChannelTest
	user := testhelper.SenderAlice

	// BroadcastSystemMessage fails (simulated by mockPresenter.Format returning nil), and RemoveUserFromRoom also fails
	{
		mockCM := realtimemock.NewMockChannelManager(ctrl)
		mockSub := realtimemock.NewMockChannelEventSubscriber(ctrl)
		mockPresenter := presentermock.NewMockMessagePresenterInterface(ctrl)
		mockMemberRepo := domainchatmock.NewMockChatMemberRepository(ctrl)
		mockLogger := logmock.NewMockLogger(ctrl)

		mockCM.EXPECT().GetOrCreateChannel(room).Times(1)
		mockSub.EXPECT().Start(room).Times(1)
		mockMemberRepo.EXPECT().AddUserToRoom(ctx, room, user).Return(nil).Times(1)
		mockMemberRepo.EXPECT().GetRoomUserCount(ctx, room).Return(int64(5), nil).Times(1)
		mockPresenter.EXPECT().Format(gomock.Any()).Return(nil).Times(1)
		mockMemberRepo.EXPECT().RemoveUserFromRoom(ctx, room, user).Return(assert.AnError).Times(1)
		mockLogger.EXPECT().Warnw("compensate RemoveUserFromRoom failed", "err", assert.AnError, "channel", room, "nickname", user).Times(1)

		service := &chatService{
			channelManager: mockCM,
			redisSub:       mockSub,
			presenter:      mockPresenter,
			memberRepo:     mockMemberRepo,
			logger:         mockLogger,
			goFunc:         func(f func()) { f() },
		}
		err := service.JoinChannel(ctx, room, user)
		assert.Error(t, err)
	}

	// All succeed
	{
		mockCM := realtimemock.NewMockChannelManager(ctrl)
		mockSub := realtimemock.NewMockChannelEventSubscriber(ctrl)
		mockPresenter := presentermock.NewMockMessagePresenterInterface(ctrl)
		mockMemberRepo := domainchatmock.NewMockChatMemberRepository(ctrl)

		mockCM.EXPECT().GetOrCreateChannel(room).Times(1)
		mockSub.EXPECT().Start(room).Times(1)
		mockMemberRepo.EXPECT().AddUserToRoom(ctx, room, user).Return(nil).Times(1)
		mockPresenter.EXPECT().Format(gomock.Any()).Return(&realtime.Message{}).Times(1)
		mockCM.EXPECT().Broadcast(gomock.Any()).Times(1)

		service := &chatService{
			channelManager: mockCM,
			redisSub:       mockSub,
			presenter:      mockPresenter,
			memberRepo:     mockMemberRepo,
			goFunc:         func(f func()) { f() },
		}
		err := service.JoinChannel(ctx, room, user)
		assert.NoError(t, err)
	}

	// AddUserToRoom failed
	{
		mockCM := realtimemock.NewMockChannelManager(ctrl)
		mockSub := realtimemock.NewMockChannelEventSubscriber(ctrl)
		mockPresenter := presentermock.NewMockMessagePresenterInterface(ctrl)
		mockMemberRepo := domainchatmock.NewMockChatMemberRepository(ctrl)

		mockCM.EXPECT().GetOrCreateChannel(room).Times(1)
		mockSub.EXPECT().Start(room).Times(1)
		mockMemberRepo.EXPECT().AddUserToRoom(ctx, room, user).Return(assert.AnError).Times(1)
		service := &chatService{
			channelManager: mockCM,
			redisSub:       mockSub,
			presenter:      mockPresenter,
			memberRepo:     mockMemberRepo,
			goFunc:         func(f func()) { f() },
		}
		err := service.JoinChannel(ctx, room, user)
		assert.Error(t, err)
	}
}

func TestChatService_BroadcastSystemMessage(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockCM := realtimemock.NewMockChannelManager(ctrl)
	mockSub := realtimemock.NewMockChannelEventSubscriber(ctrl)
	mockPresenter := presentermock.NewMockMessagePresenterInterface(ctrl)

	mockPresenter.EXPECT().Format(gomock.Any()).Return(&realtime.Message{
		Channel: testhelper.ChannelTest,
		Data: []byte("{" + testhelper.KeySender + ":" + testhelper.SenderSystem + "," +
			testhelper.KeyMessage + ":" + testhelper.SenderAlice + testhelper.MessageJoin + "}"),
	}).Times(1)

	mockCM.EXPECT().Broadcast(gomock.Any()).Times(1)

	mockMemberRepo := domainchatmock.NewMockChatMemberRepository(ctrl)
	service := &chatService{
		channelManager: mockCM,
		redisSub:       mockSub,
		presenter:      mockPresenter,
		memberRepo:     mockMemberRepo,
	}

	ctx := context.Background()
	err := service.BroadcastSystemMessage(ctx, testhelper.ChannelTest, testhelper.SenderAlice, testhelper.ActionJoin)
	assert.NoError(t, err)
}
