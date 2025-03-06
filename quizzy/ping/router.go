package ping

import (
	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
	"maps"
	"quizzy.app/backend/quizzy/services"
	"strings"
)

type Controller struct {
	FbService   *services.FirebaseServices
	RedisClient *redis.Client
}

func Configure(fbs *services.FirebaseServices, rc *redis.Client) *Controller {
	return &Controller{FbService: fbs, RedisClient: rc}
}

func (pc *Controller) ConfigureRouting(rt *gin.RouterGroup) {
	rt.GET("/ping", pc.ping)
}

func (pc *Controller) ping(c *gin.Context) {
	details := map[string]string{
		"database": "OK",
		"redis":    "OK",
	}

	// Checking firebase services (firestore + firebase) availability.
	if pc.FbService == nil {
		details["database"] = "KO"
	}

	if pc.FbService == nil {
		details["redis"] = "KO"
	}

	c.JSON(200, gin.H{
		"status":  getGlobalStatus(details),
		"details": details,
	})
}

func getGlobalStatus(details map[string]string) string {
	koc := 0

	for v := range maps.Values(details) {
		switch strings.ToUpper(v) {
		case "OK":
			continue
		case "KO":
			koc++
		}
	}

	if koc == 0 {
		return "OK"
	} else if koc == len(details) {
		return "KO"
	} else {
		return "Partial"
	}
}
