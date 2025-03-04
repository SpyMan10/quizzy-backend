package ping

import (
	"github.com/gin-gonic/gin"
	"maps"
	"strings"
)

func ConfigureRoutes(rt *gin.RouterGroup) {
	rt.GET("/ping", ping)
}

func ping(c *gin.Context) {
	details := map[string]string{
		"database": "OK",
		"redis":    "OK",
	}

	// Checking firebase services (firestore + firebase) availability.
	if _, exists := c.Get("firebase-services"); !exists {
		details["database"] = "KO"
	}

	if _, exists := c.Get("redis-service"); !exists {
		details["redis"] = "KO"
	}

	c.JSON(200, gin.H{
		"status":  compareServiceStatus(details),
		"details": details,
	})
}

func compareServiceStatus(details map[string]string) string {
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
