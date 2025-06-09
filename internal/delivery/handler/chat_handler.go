package handler

import (
	"encoding/json"
	"net/http"

	appchat "github.com/Soyuen/go-redis-chat-server/internal/application/chat"
	apperr "github.com/Soyuen/go-redis-chat-server/internal/errors"
	"github.com/Soyuen/go-redis-chat-server/pkg/loggeriface"
	"github.com/Soyuen/go-redis-chat-server/pkg/realtimeiface"
	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
)

type ChatHandler struct {
	manager     realtimeiface.ChatChannelManager // 抽象，用於邏輯用途
	connection  realtimeiface.Connection
	chatService appchat.ChatService
	upgrader    websocket.Upgrader
	logger      loggeriface.Logger
}

// 接受 Manager 介面注入
func NewChatHandler(manager realtimeiface.ChatChannelManager, connection realtimeiface.Connection,
	chatService appchat.ChatService, logger loggeriface.Logger) *ChatHandler {
	return &ChatHandler{
		manager:     manager,
		connection:  connection,
		chatService: chatService,
		logger:      logger,
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

	// 假設你之後有支援登入，可以從 session 或 token 抓 username
	sender := c.Query("nickname") // 或者從 Header/Claims 拿
	if sender == "" {
		c.JSON(http.StatusBadRequest, apperr.ErrorResponse{
			Code: apperr.ErrCodeInvalidRequestBody,
		})
		return
	}
	h.connection.HandleConnection(conn, channel, func(raw []byte) *realtimeiface.Message {
		// 呼叫 application 層進行訊息解析
		dmsg, err := h.chatService.ProcessIncoming(raw, sender, channel)
		if err != nil {
			h.logger.Warnw("failed to parse message", "err", err)
			return nil
		}

		// 👇 建立一個含 sender 的 JSON 結構
		messageObj := map[string]string{
			"sender":  sender,
			"message": dmsg.Content, // 原始內容（未 base64）或你要傳的內容
		}
		jsonBytes, err := json.Marshal(messageObj)
		if err != nil {
			h.logger.Warnw("failed to marshal message JSON", "err", err)
			return nil
		}

		// 👇 回傳給前端（sender、channel、data 全都有）
		return &realtimeiface.Message{
			Channel: dmsg.Channel,
			Data:    jsonBytes,
		}
	})

}
