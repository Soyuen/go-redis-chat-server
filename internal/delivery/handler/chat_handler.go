package handler

import (
	"net/http"

	appchat "github.com/Soyuen/go-redis-chat-server/internal/application/chat"
	apperr "github.com/Soyuen/go-redis-chat-server/internal/errors"
	"github.com/Soyuen/go-redis-chat-server/internal/presenter"
	"github.com/Soyuen/go-redis-chat-server/pkg/loggeriface"
	"github.com/Soyuen/go-redis-chat-server/pkg/realtimeiface"
	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
)

type ChatHandler struct {
	manager     realtimeiface.ChannelManager
	connection  realtimeiface.Connection
	chatService appchat.ChatService
	upgrader    websocket.Upgrader
	presenter   presenter.MessagePresenterInterface
	logger      loggeriface.Logger
}

func NewChatHandler(manager realtimeiface.ChannelManager, connection realtimeiface.Connection,
	chatService appchat.ChatService, presenter presenter.MessagePresenterInterface, logger loggeriface.Logger) *ChatHandler {
	return &ChatHandler{
		manager:     manager,
		connection:  connection,
		chatService: chatService,
		logger:      logger,
		presenter:   presenter,
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

	nickname := c.Query("nickname")
	if nickname == "" || nickname == "System" {
		c.JSON(http.StatusBadRequest, apperr.ErrorResponse{
			Code: apperr.ErrCodeInvalidRequestBody,
		})
		return
	}

	if err := h.chatService.BroadcastSystemMessage(channel, nickname, "joined"); err != nil {
		h.logger.Warnw("failed to announce join", "err", err)
	}

	h.connection.HandleConnection(conn, channel, func(raw []byte) *realtimeiface.Message {
		dmsg, err := h.chatService.ProcessIncoming(raw, nickname, channel)
		if err != nil {
			h.logger.Warnw("failed to parse message", "err", err)
			return nil
		}
		return h.presenter.Format(dmsg)
	}, func() {
		if err := h.chatService.BroadcastSystemMessage(channel, nickname, "left"); err != nil {
			h.logger.Warnw("failed to announce leave", "err", err)
		}
	})

}
