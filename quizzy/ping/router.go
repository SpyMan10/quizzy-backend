package ping

import (
	"github.com/gin-gonic/gin"
)

func ConfigureRoutes(rt *gin.RouterGroup) {
	rt.GET("/ping", ping)
}

func ping(c *gin.Context) {
	if _, exists := c.Get("firebase-services"); !exists {
		c.JSON(500, gin.H{
			"status": "Partial",
			"details": gin.H{
				"database": "KO",
			},
		})
		return
	}

	c.JSON(200, gin.H{
		"status": "OK",
		"details": gin.H{
			"database": "OK",
		},
	})
}
