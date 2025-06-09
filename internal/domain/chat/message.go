// domain/chat/message.go
package chat

import (
	"errors"
	"strings"
)

type Message struct {
	Sender  string
	Channel string
	Content string
}

func NewMessage(sender, channel, content string) (*Message, error) {
	if strings.TrimSpace(content) == "" {
		return nil, errors.New("message content cannot be empty")
	}
	return &Message{Sender: sender, Channel: channel, Content: content}, nil
}
