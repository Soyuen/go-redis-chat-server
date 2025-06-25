package realtime

import (
	"errors"
	"testing"
	"time"

	logmock "github.com/Soyuen/go-redis-chat-server/internal/infrastructure/logger/mocks"
	"github.com/Soyuen/go-redis-chat-server/internal/infrastructure/realtime/mocks"

	"github.com/Soyuen/go-redis-chat-server/pkg/realtimeiface"
	"github.com/golang/mock/gomock"
)

func setupMocks(t *testing.T) (
	ctrl *gomock.Controller,
	mockLogger *logmock.MockLogger,
	mockManager *mocks.MockChannelManager,
	mockBroadcaster *mocks.MockBroadcaster,
	mockClientFactory *mocks.MockClientFactory,
	mockClient *mocks.MockClient,
	conn *mockConn,
) {
	ctrl = gomock.NewController(t)

	mockLogger = logmock.NewMockLogger(ctrl)
	mockManager = mocks.NewMockChannelManager(ctrl)
	mockBroadcaster = mocks.NewMockBroadcaster(ctrl)
	mockClientFactory = mocks.NewMockClientFactory(ctrl)
	mockClient = mocks.NewMockClient(ctrl)
	conn = &mockConn{writeChan: make(chan []byte, 1)}

	return
}

func waitWritePumpCalled(t *testing.T, writePumpCalled chan struct{}) {
	select {
	case <-writePumpCalled:
	case <-time.After(time.Second):
		t.Fatal("timeout waiting for WritePump to be called")
	}
}

func TestConnection_HandleConnection_RegistersClientAndStartsPumps(t *testing.T) {
	ctrl, mockLogger, mockManager, mockBroadcaster, mockClientFactory, mockClient, conn := setupMocks(t)
	defer ctrl.Finish()

	channel := "room1"

	mockClientFactory.EXPECT().New(conn).Return(mockClient)
	mockManager.EXPECT().GetOrCreateChannel(channel).Return(mockBroadcaster)
	mockBroadcaster.EXPECT().Register(mockClient)

	writePumpCalled := make(chan struct{})
	mockClient.EXPECT().WritePump().Do(func() {
		close(writePumpCalled)
	})

	mockClient.EXPECT().ReadPump(gomock.Any()).Return(nil)
	mockBroadcaster.EXPECT().Unregister(mockClient)
	mockClient.EXPECT().Close()

	c := NewConnection(mockManager, mockLogger, mockClientFactory)

	go c.HandleConnection(conn, channel, func(raw []byte) *realtimeiface.Message {
		return nil
	}, nil)

	waitWritePumpCalled(t, writePumpCalled)
}

func TestConnection_HandleConnection_BroadcastsOnMessage(t *testing.T) {
	ctrl, mockLogger, mockManager, mockBroadcaster, mockClientFactory, mockClient, conn := setupMocks(t)
	defer ctrl.Finish()

	channel := "room1"
	rawMsg := []byte(`{"channel":"room1","data":"hi"}`)

	writePumpCalled := make(chan struct{})
	unregisterCalled := make(chan struct{})

	mockClientFactory.EXPECT().New(conn).Return(mockClient)
	mockManager.EXPECT().GetOrCreateChannel(channel).Return(mockBroadcaster)
	mockBroadcaster.EXPECT().Register(mockClient)

	mockClient.EXPECT().WritePump().Do(func() {
		close(writePumpCalled)
	})

	mockClient.EXPECT().ReadPump(gomock.Any()).DoAndReturn(func(cb func([]byte)) error {
		cb(rawMsg)
		return nil
	})

	mockManager.EXPECT().Broadcast(realtimeiface.Message{
		Channel: "room1",
		Data:    rawMsg,
	})

	mockBroadcaster.EXPECT().Unregister(mockClient).Do(func(_ realtimeiface.Client) {
		close(unregisterCalled)
	})
	mockClient.EXPECT().Close()

	c := NewConnection(mockManager, mockLogger, mockClientFactory)
	go c.HandleConnection(conn, channel, func(raw []byte) *realtimeiface.Message {
		return &realtimeiface.Message{Channel: "room1", Data: raw}
	}, nil)

	waitWritePumpCalled(t, writePumpCalled)

	select {
	case <-unregisterCalled:
	case <-time.After(time.Second):
		t.Fatal("timeout waiting for Unregister to be called")
	}
}

func TestConnection_HandleConnection_UnregistersAndClosesOnError(t *testing.T) {
	ctrl, mockLogger, mockManager, mockBroadcaster, mockClientFactory, mockClient, conn := setupMocks(t)
	defer ctrl.Finish()

	channel := "room1"

	mockClientFactory.EXPECT().New(conn).Return(mockClient)
	mockManager.EXPECT().GetOrCreateChannel(channel).Return(mockBroadcaster)
	mockBroadcaster.EXPECT().Register(mockClient)

	writePumpCalled := make(chan struct{})
	mockClient.EXPECT().WritePump().Do(func() {
		close(writePumpCalled)
	})

	mockClient.EXPECT().ReadPump(gomock.Any()).Return(errors.New("read error"))
	mockBroadcaster.EXPECT().Unregister(mockClient)
	mockClient.EXPECT().Close()

	c := NewConnection(mockManager, mockLogger, mockClientFactory)

	go c.HandleConnection(conn, channel, func(raw []byte) *realtimeiface.Message {
		return nil
	}, nil)

	// Wait for WritePump to be called to prevent the test from finishing too early
	// before the goroutine starts.
	select {
	case <-writePumpCalled:
	case <-time.After(time.Second):
		t.Fatal("timeout waiting for WritePump to be called")
	}
}

func TestConnection_HandleConnection_CallsOnClose(t *testing.T) {
	ctrl, mockLogger, mockManager, mockBroadcaster, mockClientFactory, mockClient, conn := setupMocks(t)
	defer ctrl.Finish()

	channel := "room1"

	mockClientFactory.EXPECT().New(conn).Return(mockClient)
	mockManager.EXPECT().GetOrCreateChannel(channel).Return(mockBroadcaster)
	mockBroadcaster.EXPECT().Register(mockClient)

	writePumpCalled := make(chan struct{})
	mockClient.EXPECT().WritePump().Do(func() {
		close(writePumpCalled)
	})

	mockClient.EXPECT().ReadPump(gomock.Any()).Return(nil)
	mockBroadcaster.EXPECT().Unregister(mockClient)
	mockClient.EXPECT().Close()

	onCloseCalled := make(chan struct{})
	onClose := func() {
		close(onCloseCalled)
	}

	c := NewConnection(mockManager, mockLogger, mockClientFactory)
	go c.HandleConnection(conn, channel, func(raw []byte) *realtimeiface.Message {
		return nil
	}, onClose)

	select {
	case <-writePumpCalled:
	case <-time.After(time.Second):
		t.Fatal("timeout waiting for WritePump to be called")
	}

	select {
	case <-onCloseCalled:
	case <-time.After(time.Second):
		t.Fatal("timeout waiting for onClose to be called")
	}
}

func TestConnection_handleWrite_RecoversFromPanic(t *testing.T) {
	ctrl, mockLogger, mockManager, _, mockClientFactory, mockClient, _ := setupMocks(t)
	defer ctrl.Finish()

	c := NewConnection(mockManager, mockLogger, mockClientFactory)

	// Simulate WritePump panicking
	mockClient.EXPECT().WritePump().Do(func() {
		panic("test panic")
	})

	// Expect Errorw to be called once with msg matching the expected string
	mockLogger.EXPECT().Errorw("Recovered from panic in WritePump", gomock.Any()).Times(1)

	// Call handleWrite to trigger panic and recover
	c.handleWrite(mockClient)
}
