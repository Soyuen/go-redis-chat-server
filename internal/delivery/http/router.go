package delivery

import (
	"github.com/Soyuen/go-redis-chat-server/pkg/cache"
	"github.com/gin-gonic/gin"
)

func NewRouter(redis cache.RedisCache) *gin.Engine {
	gin.SetMode(gin.DebugMode)
	r := gin.Default()

	// 注入 redis 給路由註冊用
	registerChatRoutes(r, redis)

	return r
}
func registerChatRoutes(r *gin.Engine, redis cache.RedisCache) {
}
