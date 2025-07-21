package router

import (
	"github.com/Soyuen/go-redis-chat-server/internal/application/chat"
	"github.com/Soyuen/go-redis-chat-server/internal/delivery/handler"
	"github.com/Soyuen/go-redis-chat-server/internal/presenter"

	"github.com/Soyuen/go-redis-chat-server/internal/application/realtime"
	"github.com/Soyuen/go-redis-chat-server/pkg/loggeriface"
	"github.com/gin-gonic/gin"
)

func NewRouter(manager realtime.ChannelManager,
	connection realtime.Connection,
	chatService chat.ChatService, presenter presenter.MessagePresenterInterface, logger loggeriface.Logger) *gin.Engine {
	gin.SetMode(gin.DebugMode)
	r := gin.Default()

	registerChatRoutes(r, manager, connection, chatService, presenter, logger)

	return r
}
func registerChatRoutes(r *gin.Engine, manager realtime.ChannelManager,
	connection realtime.Connection,
	chatService chat.ChatService, presenter presenter.MessagePresenterInterface, logger loggeriface.Logger) {
	chatHandler := handler.NewChatHandler(manager, connection, chatService, presenter, logger)
	chatGroup := r.Group("/chat")
	chatGroup.GET("/join", chatHandler.JoinChannel)
}
