package realtime

import (
	"errors"
	"sync"
	"testing"
	"time"

	logmock "github.com/Soyuen/go-redis-chat-server/internal/infrastructure/logger/mocks"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
)

// mock websocket.Conn 用 minimal 介面
type mockConn struct {
	writeChan chan []byte
	closeOnce sync.Once
	closed    bool
}

func (m *mockConn) ReadMessage() (int, []byte, error) {
	return 0, nil, errors.New("not implemented")
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

	// Send 應該成功寫入 channel
	client.Send([]byte("hello"))
	assert.Equal(t, 1, len(client.send))

	// channel 滿了，Send 要關閉 client
	client.send = make(chan []byte, 1)
	client.Send([]byte("msg1"))
	client.Send([]byte("msg2")) // 第二次會觸發 Close()

	// Close 再呼叫不會 panic
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

	// 寫入一條訊息
	client.send <- []byte("test message")

	// 啟動 WritePump (跑在 goroutine)
	go client.WritePump()

	// 從 mockConn 的 writeChan 收到訊息
	select {
	case msg := <-conn.writeChan:
		assert.Equal(t, "test message", string(msg))
	case <-time.After(time.Second):
		t.Fatal("timeout waiting for WritePump to write message")
	}
}
