package chat

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewMessage(t *testing.T) {
	// Test normal case
	msg, err := NewMessage("alice", "room1", "hello world")
	assert.NoError(t, err)
	assert.NotNil(t, msg)
	assert.Equal(t, "alice", msg.Sender)
	assert.Equal(t, "room1", msg.Channel)
	assert.Equal(t, "hello world", msg.Content)

	// Test empty content returns error
	msg, err = NewMessage("bob", "room1", "   ")
	assert.Error(t, err)
	assert.Nil(t, msg)
}
