package ping

import (
	"github.com/gin-gonic/gin"
)

func Setup(rt *gin.Engine) {
	rt.GET("/ping", ping)
}

func ping(c *gin.Context) {
	c.JSON(500, gin.H{
		"status": "OK",
	})
}
