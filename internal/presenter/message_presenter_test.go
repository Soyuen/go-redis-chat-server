package presenter

import (
	"encoding/json"
	"testing"

	logmock "github.com/Soyuen/go-redis-chat-server/internal/infrastructure/logger/mocks"
	"github.com/Soyuen/go-redis-chat-server/internal/testhelper"
	"github.com/golang/mock/gomock"

	"github.com/Soyuen/go-redis-chat-server/internal/domain/chat"
	"github.com/stretchr/testify/assert"
)

func TestMessagePresenter_Format(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockLogger := logmock.NewMockLogger(ctrl)

	mockLogger.EXPECT().Warnw(gomock.Any(), gomock.Any()).AnyTimes()
	mockLogger.EXPECT().Errorw(gomock.Any(), gomock.Any()).AnyTimes()
	mockLogger.EXPECT().Infow(gomock.Any(), gomock.Any()).AnyTimes()

	presenter := NewMessagePresenter(mockLogger)

	chatMsg := &chat.Message{
		Sender:  testhelper.SenderAlice,
		Channel: testhelper.ChannelTest,
		Content: testhelper.MessageTest,
	}

	result := presenter.Format(chatMsg)
	assert.NotNil(t, result)
	assert.Equal(t, chatMsg.Channel, result.Channel)

	var parsed map[string]string
	err := json.Unmarshal(result.Data, &parsed)
	assert.NoError(t, err)
	assert.Equal(t, testhelper.SenderAlice, parsed[testhelper.KeySender])
	assert.Equal(t, testhelper.MessageTest, parsed[testhelper.KeyMessage])
}
