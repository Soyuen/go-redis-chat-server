package realtime

import (
	"errors"
	"io"
	"sync"
	"testing"
	"time"

	logmock "github.com/Soyuen/go-redis-chat-server/internal/infrastructure/logger/mocks"
	"github.com/golang/mock/gomock"
	"github.com/gorilla/websocket"
	"github.com/stretchr/testify/assert"
)

// mock websocket.Conn with minimal interface
type mockConn struct {
	writeChan       chan []byte
	readMessageFunc func() (int, []byte, error)
	closeOnce       sync.Once
	closed          bool
}

func (m *mockConn) ReadMessage() (int, []byte, error) {
	return m.readMessageFunc()
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

func TestClient_Send_Close(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockLogger := logmock.NewMockLogger(ctrl)

	conn := &mockConn{writeChan: make(chan []byte, 1)}

	client := &Client{
		conn:   conn,
		logger: mockLogger,
		send:   make(chan []byte, 1),
	}

	// Send should successfully write to the channel
	client.Send([]byte("hello"))
	assert.Equal(t, 1, len(client.send))

	// When the channel is full, Send should close the client
	client.send = make(chan []byte, 1)
	client.Send([]byte("msg1"))
	client.Send([]byte("msg2")) // The second call should trigger Close()

	// Calling Close multiple times should not panic
	client.Close()
	client.Close()
}

func TestClient_WritePump(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockLogger := logmock.NewMockLogger(ctrl)
	conn := &mockConn{writeChan: make(chan []byte, 1)}

	client := &Client{
		conn:   conn,
		logger: mockLogger,
		send:   make(chan []byte, 1),
	}

	// Write a single message
	client.send <- []byte("test message")

	// Start WritePump (runs in a goroutine)
	go client.WritePump()

	// Receive message from mockConn's writeChan
	select {
	case msg := <-conn.writeChan:
		assert.Equal(t, "test message", string(msg))
	case <-time.After(time.Second):
		t.Fatal("timeout waiting for WritePump to write message")
	}
}

// Verify that onMessage is called when a message is received.
func TestClient_ReadPump_TriggersOnMessage(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockLogger := logmock.NewMockLogger(ctrl)
	// Allow Errorw and Infow methods to be called without failing
	mockLogger.EXPECT().Errorw(gomock.Any(), gomock.Any()).AnyTimes()
	mockLogger.EXPECT().Infow(gomock.Any(), gomock.Any()).AnyTimes()

	calls := 0
	conn := &mockConn{
		readMessageFunc: func() (int, []byte, error) {
			if calls == 0 {
				calls++
				return websocket.TextMessage, []byte("hello"), nil
			}
			return 0, nil, io.EOF // ReadPump exits upon receiving EOF.
		},
	}

	client := &Client{
		conn:   conn,
		logger: mockLogger,
		send:   make(chan []byte, 1),
	}

	messageHandled := make(chan struct{})

	go func() {
		err := client.ReadPump(func(data []byte) {
			assert.Equal(t, "hello", string(data))
			close(messageHandled)
		})
		// ReadPump is expected to exit due to io.EOF
		assert.ErrorIs(t, err, io.EOF)
	}()

	select {
	case <-messageHandled:
		// success
	case <-time.After(time.Second):
		t.Fatal("timeout: onMessage not triggered")
	}
}
