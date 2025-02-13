package http

import (
	"github.com/gin-gonic/gin"
	"quizzy.app/backend/quizzy/ping"
	"quizzy.app/backend/quizzy/quiz"
	"quizzy.app/backend/quizzy/users"
)

func ConfigureRouting(router *gin.RouterGroup) {
	grp := router.Group("/api")
	ping.ConfigureRoutes(grp)
	users.ConfigureRoutes(grp)
	quiz.ConfigureRoutes(grp)
}
