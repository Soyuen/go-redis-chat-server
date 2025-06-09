package router

import (
	"github.com/Soyuen/go-redis-chat-server/internal/application/chat"
	"github.com/Soyuen/go-redis-chat-server/internal/delivery/handler"

	"github.com/Soyuen/go-redis-chat-server/pkg/loggeriface"
	"github.com/Soyuen/go-redis-chat-server/pkg/realtimeiface"
	"github.com/gin-gonic/gin"
)

func NewRouter(manager realtimeiface.ChatChannelManager,
	connection realtimeiface.Connection,
	chatService chat.ChatService, logger loggeriface.Logger) *gin.Engine {
	gin.SetMode(gin.DebugMode)
	r := gin.Default()

	registerChatRoutes(r, manager, connection, chatService, logger)

	return r
}
func registerChatRoutes(r *gin.Engine, manager realtimeiface.ChatChannelManager,
	connection realtimeiface.Connection,
	chatService chat.ChatService, logger loggeriface.Logger) {
	chatHandler := handler.NewChatHandler(manager, connection, chatService, logger)
	chatGroup := r.Group("/chat")
	chatGroup.GET("/join", chatHandler.JoinChannel)
}
