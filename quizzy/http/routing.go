package http

import (
	"github.com/gin-gonic/gin"
	"quizzy.app/backend/quizzy/ping"
)

func ConfigureRouting(router *gin.RouterGroup) {
	ping.ConfigureRoutes(router)
}
