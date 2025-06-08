package handler

import (
	"net/http"

	"github.com/Soyuen/go-redis-chat-server/internal/application/chat"
	apperr "github.com/Soyuen/go-redis-chat-server/internal/errors"
	"github.com/Soyuen/go-redis-chat-server/pkg/realtimeiface"
	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
)

type ChatHandler struct {
	manager     realtimeiface.ChatChannelManager // 抽象，用於邏輯用途
	connection  realtimeiface.Connection
	chatService chat.ChatService
	upgrader    websocket.Upgrader
}

// 接受 Manager 介面注入
func NewChatHandler(manager realtimeiface.ChatChannelManager, connection realtimeiface.Connection,
	chatService chat.ChatService) *ChatHandler {
	return &ChatHandler{
		manager:     manager,
		connection:  connection,
		chatService: chatService,
		upgrader: websocket.Upgrader{
			ReadBufferSize:  1024,
			WriteBufferSize: 1024,
			CheckOrigin: func(r *http.Request) bool {
				// TODO allow only specific origins in the future, for example by checking r.Header.Get("Origin").
				return true
			},
		},
	}
}

func (h *ChatHandler) JoinChannel(c *gin.Context) {
	channel := c.Query("channel")
	if channel == "" {
		c.JSON(http.StatusBadRequest, apperr.ErrorResponse{
			Code: apperr.ErrCodeInvalidRequestBody,
		})
		return
	}

	conn, err := h.upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		c.JSON(http.StatusInternalServerError, apperr.ErrorResponse{
			Code: apperr.ErrCodeWebSocketUpgradeFailed,
		})
		return
	}
	if err := h.chatService.CreateRoom(channel); err != nil {
		conn.Close()
		c.JSON(http.StatusInternalServerError, apperr.ErrorResponse{
			Code: apperr.ErrCodeChannelCreationFailed,
		})
		return
	}

	h.connection.HandleConnection(conn, channel)
}
