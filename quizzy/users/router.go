package ping

import (
	"github.com/gin-gonic/gin"
)

type User struct {
}

func ConfigureRoutes(rt *gin.RouterGroup) {
	rt.POST("/users", createUser)
}

func createUser(c *gin.Context) {
	data := map[string]interface{}{}
	_ = c.ShouldBindJSON(&data)
	println(data)
}
