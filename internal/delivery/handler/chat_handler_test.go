package handler

import (
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	chatmock "github.com/Soyuen/go-redis-chat-server/internal/application/chat/mocks"
	"github.com/Soyuen/go-redis-chat-server/internal/application/realtime"
	realtimemock "github.com/Soyuen/go-redis-chat-server/internal/application/realtime/mocks"
	domainchat "github.com/Soyuen/go-redis-chat-server/internal/domain/chat"
	logmock "github.com/Soyuen/go-redis-chat-server/internal/infrastructure/logger/mocks"
	presentermock "github.com/Soyuen/go-redis-chat-server/internal/presenter/mocks"
	"github.com/Soyuen/go-redis-chat-server/internal/testhelper"

	"github.com/gin-gonic/gin"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
)

func TestChatHandler_JoinChannel_CreateRoomError(t *testing.T) {
	gin.SetMode(gin.TestMode)

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockChatService := chatmock.NewMockChatService(ctrl)
	mockChatService.EXPECT().BroadcastSystemMessage(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).AnyTimes()
	mockConnection := realtimemock.NewMockConnection(ctrl)
	mockPresenter := presentermock.NewMockMessagePresenterInterface(ctrl)
	mockLogger := logmock.NewMockLogger(ctrl)
	mockWSConn := realtimemock.NewMockWSConn(ctrl)

	handlerChat := NewChatHandler(nil, mockConnection, mockChatService, mockPresenter, mockLogger)

	// Inject a mock upgraderFunc
	handlerChat.SetUpgraderFunc(func(w http.ResponseWriter, r *http.Request) (realtime.WSConn, error) {
		return mockWSConn, nil
	})

	// Simulate failure when creating a room
	mockChatService.EXPECT().CreateRoom("room1").Return(errors.New("fail")).Times(1)

	// Expect conn.Close() to be called
	mockWSConn.EXPECT().Close().Return(nil).Times(1)

	router := gin.New()
	router.GET("/chat/join", handlerChat.JoinChannel)

	req := httptest.NewRequest(http.MethodGet, "/chat/join?channel=room1&nickname=Alice", nil)
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusInternalServerError, w.Code)
}

func TestChatHandler_JoinChannel_EmptyChannel(t *testing.T) {
	gin.SetMode(gin.TestMode)

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockChatService := chatmock.NewMockChatService(ctrl)
	mockConnection := realtimemock.NewMockConnection(ctrl)
	mockPresenter := presentermock.NewMockMessagePresenterInterface(ctrl)
	mockLogger := logmock.NewMockLogger(ctrl)

	handlerChat := NewChatHandler(nil, mockConnection, mockChatService, mockPresenter, mockLogger)

	router := gin.New()
	router.GET("/chat/join", handlerChat.JoinChannel)

	req := httptest.NewRequest(http.MethodGet, "/chat/join?channel=&nickname=Alice", nil) // channel空字串
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestChatHandler_JoinChannel_UpgradeFail(t *testing.T) {
	gin.SetMode(gin.TestMode)

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockChatService := chatmock.NewMockChatService(ctrl)
	mockConnection := realtimemock.NewMockConnection(ctrl)
	mockPresenter := presentermock.NewMockMessagePresenterInterface(ctrl)
	mockLogger := logmock.NewMockLogger(ctrl)

	handlerChat := NewChatHandler(nil, mockConnection, mockChatService, mockPresenter, mockLogger)

	// Mock upgraderFunc to return an error, simulating a failed WebSocket upgrade
	handlerChat.SetUpgraderFunc(func(w http.ResponseWriter, r *http.Request) (realtime.WSConn, error) {
		return nil, errors.New("upgrade failed")
	})

	router := gin.New()
	router.GET("/chat/join", handlerChat.JoinChannel)

	req := httptest.NewRequest(http.MethodGet, "/chat/join?channel=room1&nickname=Alice", nil)
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusInternalServerError, w.Code)
}

func TestChatHandler_JoinChannel_InvalidNickname(t *testing.T) {
	gin.SetMode(gin.TestMode)
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockChatService := chatmock.NewMockChatService(ctrl)
	mockConnection := realtimemock.NewMockConnection(ctrl)
	mockPresenter := presentermock.NewMockMessagePresenterInterface(ctrl)
	mockLogger := logmock.NewMockLogger(ctrl)

	handlerChat := NewChatHandler(nil, mockConnection, mockChatService, mockPresenter, mockLogger)

	// Use a mock WSConn to avoid opening a real WebSocket connection
	mockWsConn := realtimemock.NewMockWSConn(ctrl)
	handlerChat.SetUpgraderFunc(func(w http.ResponseWriter, r *http.Request) (realtime.WSConn, error) {
		return mockWsConn, nil
	})

	router := gin.New()
	router.GET("/chat/join", handlerChat.JoinChannel)

	// 1. nickname is empty
	req := httptest.NewRequest(http.MethodGet, "/chat/join?channel=room1&nickname=", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)
	assert.Equal(t, http.StatusBadRequest, w.Code)

	// 2. nickname is "System"
	req2 := httptest.NewRequest(http.MethodGet, "/chat/join?channel=room1&nickname=System", nil)
	w2 := httptest.NewRecorder()
	router.ServeHTTP(w2, req2)
	assert.Equal(t, http.StatusBadRequest, w2.Code)
}

func TestChatHandler_MessageHandler(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockChatService := chatmock.NewMockChatService(ctrl)
	mockPresenter := presentermock.NewMockMessagePresenterInterface(ctrl)
	mockLogger := logmock.NewMockLogger(ctrl)

	handlerChat := NewChatHandler(nil, nil, mockChatService, mockPresenter, mockLogger)

	channel := testhelper.ChannelTest
	nickname := testhelper.SenderAlice
	raw := []byte(testhelper.MessageHello)

	// successful
	mockChatService.EXPECT().
		ProcessIncoming(gomock.Any(), gomock.Any(), gomock.Any()).
		Return(&domainchat.Message{
			Sender:  testhelper.SenderAlice,
			Channel: testhelper.ChannelTest,
			Content: testhelper.MessageHello,
		}, nil)
	mockPresenter.EXPECT().Format(gomock.Any()).Return(&realtime.Message{})

	f := handlerChat.messageHandler(channel, nickname)
	msg := f(raw)
	assert.NotNil(t, msg)

	// failed
	mockChatService.EXPECT().ProcessIncoming(raw, nickname, channel).Return(nil, errors.New("fail"))
	mockLogger.EXPECT().Warnw("failed to parse message", gomock.Any())

	msg = f(raw)
	assert.Nil(t, msg)
}

func TestChatHandler_LeaveHandler(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockChatService := chatmock.NewMockChatService(ctrl)
	mockLogger := logmock.NewMockLogger(ctrl)

	handlerChat := NewChatHandler(nil, nil, mockChatService, nil, mockLogger)

	channel := "room1"
	nickname := "Alice"

	// successful
	mockChatService.EXPECT().BroadcastSystemMessage(gomock.Any(), channel, nickname, "left").Return(nil)
	h := handlerChat.leaveHandler(channel, nickname)
	h()

	// failed
	mockChatService.EXPECT().BroadcastSystemMessage(gomock.Any(), channel, nickname, "left").Return(errors.New("fail"))
	mockLogger.EXPECT().Warnw("failed to announce leave", gomock.Any())
	h()
}
