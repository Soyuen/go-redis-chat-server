package realtime

import (
	"errors"
	"io"
	"sync"
	"testing"
	"time"

	logmock "github.com/Soyuen/go-redis-chat-server/internal/infrastructure/logger/mocks"
	"github.com/Soyuen/go-redis-chat-server/internal/testhelper"
	"github.com/golang/mock/gomock"
	"github.com/gorilla/websocket"
	"github.com/stretchr/testify/assert"
)

// Improved mockConn to support simulating multiple messages
type mockConn struct {
	writeChan       chan []byte
	messages        [][]byte
	readIndex       int
	closed          bool
	closeOnce       sync.Once
	readMessageFunc func() (int, []byte, error)
}

func (m *mockConn) ReadMessage() (int, []byte, error) {
	if m.readMessageFunc != nil {
		return m.readMessageFunc()
	}
	if m.readIndex >= len(m.messages) {
		return 0, nil, io.EOF
	}
	msg := m.messages[m.readIndex]
	m.readIndex++
	return websocket.TextMessage, msg, nil
}

func (m *mockConn) WriteMessage(messageType int, data []byte) error {
	if m.closed {
		return errors.New("connection closed")
	}
	m.writeChan <- data
	return nil
}

func (m *mockConn) Close() error {
	m.closeOnce.Do(func() {
		m.closed = true
		close(m.writeChan)
	})
	return nil
}

// Test Send successfully writes to the channel
func TestClient_Send_WritesToChannel(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	client := &Client{
		conn:   &mockConn{writeChan: make(chan []byte, 1)},
		send:   make(chan []byte, 1),
		logger: logmock.NewMockLogger(ctrl),
	}

	client.Send([]byte(testhelper.MessageTest))
	assert.Equal(t, 1, len(client.send))
}

// Test Send closes the client when the channel is full
func TestClient_Send_ClosesClientWhenChannelFull(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockLogger := logmock.NewMockLogger(ctrl)
	conn := &mockConn{writeChan: make(chan []byte, 1)}

	client := &Client{
		conn:   conn,
		send:   make(chan []byte, 1),
		logger: mockLogger,
	}

	client.Send([]byte(testhelper.MessageTest))
	client.Send([]byte(testhelper.MessageHello)) // Should close the client when channel is full
	assert.True(t, conn.closed)
}

// Test Close is safe to call multiple times
func TestClient_Close_Idempotent(t *testing.T) {
	client := &Client{
		conn:   &mockConn{writeChan: make(chan []byte, 1)},
		send:   make(chan []byte, 1),
		logger: logmock.NewMockLogger(gomock.NewController(t)),
	}

	client.Close()
	client.Close()
}

// Test WritePump successfully writes to the connection
func TestClient_WritePump_WritesMessage(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	conn := &mockConn{writeChan: make(chan []byte, 1)}
	client := &Client{
		conn:   conn,
		send:   make(chan []byte, 1),
		logger: logmock.NewMockLogger(ctrl),
	}

	client.send <- []byte(testhelper.MessageTest)
	go client.WritePump()

	select {
	case msg := <-conn.writeChan:
		assert.Equal(t, testhelper.MessageTest, string(msg))
	case <-time.After(time.Second):
		t.Fatal("timeout waiting for WritePump to write message")
	}
}

// Test ReadPump calls onMessage when a message is received
func TestClient_ReadPump_TriggersOnMessage(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockLogger := logmock.NewMockLogger(ctrl)
	mockLogger.EXPECT().Errorw(gomock.Any(), gomock.Any()).AnyTimes()
	mockLogger.EXPECT().Infow(gomock.Any(), gomock.Any()).AnyTimes()

	conn := &mockConn{
		messages:  [][]byte{[]byte(testhelper.MessageTest)},
		writeChan: make(chan []byte, 1),
	}

	client := &Client{
		conn:   conn,
		send:   make(chan []byte, 1),
		logger: mockLogger,
	}

	done := make(chan struct{})
	go func() {
		err := client.ReadPump(func(data []byte) {
			assert.Equal(t, testhelper.MessageTest, string(data))
			close(done)
		})
		assert.ErrorIs(t, err, io.EOF)
	}()

	select {
	case <-done:
	case <-time.After(time.Second):
		t.Fatal("onMessage not triggered")
	}
}

func TestClient_ReadPump_ReturnsErrorOnUnexpectedFailure(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockLogger := logmock.NewMockLogger(ctrl)
	mockLogger.EXPECT().Errorw(gomock.Any(), gomock.Any()).AnyTimes()
	mockLogger.EXPECT().Infow(gomock.Any(), gomock.Any()).AnyTimes()

	called := false
	conn := &mockConn{
		readMessageFunc: func() (int, []byte, error) {
			if !called {
				called = true
				return 0, nil, errors.New("unexpected error")
			}
			return 0, nil, io.EOF
		},
		writeChan: make(chan []byte, 1),
	}

	client := &Client{
		conn:   conn,
		send:   make(chan []byte, 1),
		logger: mockLogger,
	}

	err := client.ReadPump(func([]byte) {
		t.Fatal("onMessage should not be triggered on error")
	})

	assert.Error(t, err)
	assert.EqualError(t, err, "unexpected error")
}
