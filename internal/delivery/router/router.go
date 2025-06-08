package router

import (
	"github.com/Soyuen/go-redis-chat-server/internal/application/chat"
	"github.com/Soyuen/go-redis-chat-server/internal/delivery/handler"

	"github.com/Soyuen/go-redis-chat-server/pkg/realtimeiface"
	"github.com/gin-gonic/gin"
)

func NewRouter(manager realtimeiface.ChatChannelManager,
	connection realtimeiface.Connection,
	chatService chat.ChatService) *gin.Engine {
	gin.SetMode(gin.DebugMode)
	r := gin.Default()

	registerChatRoutes(r, manager, connection, chatService)

	return r
}
func registerChatRoutes(r *gin.Engine, manager realtimeiface.ChatChannelManager,
	connection realtimeiface.Connection,
	chatService chat.ChatService,
) {
	chatHandler := handler.NewChatHandler(manager, connection, chatService)
	chatGroup := r.Group("/chat")
	chatGroup.GET("/join", chatHandler.JoinChannel)
}
