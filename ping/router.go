package ping

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func Setup(rt *gin.Engine) {
	rt.GET("/ping", func(c *gin.Context) {
		if http.StatusOK == 200 {
			c.JSON(200, gin.H{
				"status": "OK",
			})
		}
	})
}
