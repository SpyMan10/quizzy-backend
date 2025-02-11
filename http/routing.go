package http

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

func Setup(rt *gin.Engine) {
	rt.GET("/ping", func(c *gin.Context) {
		c.String(http.StatusOK, "pong!")
	})
}
