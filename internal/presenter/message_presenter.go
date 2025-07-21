// presenter/message_presenter.go
package presenter

import (
	"encoding/json"

	"github.com/Soyuen/go-redis-chat-server/internal/application/realtime"
	"github.com/Soyuen/go-redis-chat-server/internal/domain/chat"
	"github.com/Soyuen/go-redis-chat-server/pkg/loggeriface"
)

type MessagePresenterInterface interface {
	Format(msg *chat.Message) *realtime.Message
}

type MessagePresenter struct {
	logger loggeriface.Logger
}

func NewMessagePresenter(logger loggeriface.Logger) MessagePresenterInterface {
	return &MessagePresenter{logger: logger}
}

func (p *MessagePresenter) Format(msg *chat.Message) *realtime.Message {
	messageObj := map[string]string{
		"sender":  msg.Sender,
		"message": msg.Content,
	}
	jsonBytes, err := json.Marshal(messageObj)
	if err != nil {
		p.logger.Warnw("failed to marshal message JSON", "err", err)
		return nil
	}
	return &realtime.Message{
		Channel: msg.Channel,
		Data:    jsonBytes,
	}
}
