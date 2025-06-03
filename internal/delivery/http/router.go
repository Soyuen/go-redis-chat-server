package delivery

import (
	"github.com/gin-gonic/gin"
)

func NewRouter() *gin.Engine {
	gin.SetMode(gin.DebugMode) // Set Gin mode: debug / release / test
	r := gin.Default()
	// You can add Swagger UI for API documentation here.
	registerChatRoutes(r)

	return r
}

func registerChatRoutes(r *gin.Engine) {
}
