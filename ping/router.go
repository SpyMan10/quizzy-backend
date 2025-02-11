package ping

import (
	"github.com/gin-gonic/gin"
)

func Setup(rt *gin.Engine) {
	rt.GET("/ping", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"status": "OK",
		})
	})
}
