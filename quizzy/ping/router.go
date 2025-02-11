package ping

import (
	"github.com/gin-gonic/gin"
)

func ConfigureRoutes(rt *gin.RouterGroup) {
	rt.GET("/ping", ping)
}

func ping(c *gin.Context) {
	c.JSON(500, gin.H{
		"status": "OK",
	})
}
