package http

import (
	"github.com/gin-gonic/gin"
	"quizzy.app/backend/ping"
)

func Setup(rt *gin.Engine) {
	ping.Setup(rt)
}
