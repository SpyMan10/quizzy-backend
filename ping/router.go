package ping

import (
	"github.com/gin-gonic/gin"
)

func Setup(rt *gin.Engine) {
	rt.GET("/ping", ping)
}

func ping(c *gin.Context) {
	// Teste de la reponse FireBAse
	payload := make(map[string]interface{})

	dataStatus := false // A remplacer par un boolean de test de connexion a la base de donnée
	statusHtml := 200
	printHtml := "OK"
	printDatabase := "OK"

	if !dataStatus {
		statusHtml = 500
		printHtml = "Partial"
		printDatabase = "KO"
	}

	payload["status"] = printHtml
	payload["details"] = map[string]string{"database": printDatabase}

	c.JSON(statusHtml, payload)
}
