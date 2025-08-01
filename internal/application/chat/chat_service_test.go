package chat

import (
	"encoding/json"
	"testing"

	realtimemock "github.com/Soyuen/go-redis-chat-server/internal/application/realtime/mocks"
	domainchatmock "github.com/Soyuen/go-redis-chat-server/internal/domain/chat/mocks"
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
	mockCM.EXPECT().GetOrCreateChannel(testhelper.ChannelTest).Times(1)
	mockSub.EXPECT().Start(testhelper.ChannelTest).Times(1)

	service := NewChatService(mockCM, mockSub, mockPresenter, mockMemberRepo)
	service.(*chatService).goFunc = func(f func()) { f() }

	err := service.CreateRoom(testhelper.ChannelTest)
	assert.NoError(t, err)
}

func TestChatService_ProcessIncoming_Valid(t *testing.T) {
	service := NewChatService(nil, nil, presentermock.NewMockMessagePresenterInterface(nil), nil)

	raw, _ := json.Marshal(map[string]string{testhelper.KeyMessage: testhelper.MessageHello})
	msg, err := service.ProcessIncoming(raw, testhelper.SenderAlice, testhelper.ChannelTest)

	assert.NoError(t, err)
	assert.Equal(t, testhelper.SenderAlice, msg.Sender)
	assert.Equal(t, testhelper.ChannelTest, msg.Channel)
	assert.Equal(t, testhelper.MessageHello, msg.Content)
}

func TestChatService_ProcessIncoming_InvalidJSON(t *testing.T) {
	service := NewChatService(nil, nil, presentermock.NewMockMessagePresenterInterface(nil), nil)

	_, err := service.ProcessIncoming([]byte("invalid json"), testhelper.SenderAlice, testhelper.ChannelTest)

	assert.EqualError(t, err, "invalid message format")
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

	err := service.BroadcastSystemMessage(testhelper.ChannelTest, testhelper.SenderAlice, testhelper.ActionJoin)
	assert.NoError(t, err)
}
