package ping

import (
	"github.com/gin-gonic/gin"
)

func ConfigureRoutes(rt *gin.RouterGroup) {
	rt.GET("/ping", ping)
}

func ping(c *gin.Context) {
<<<<<<< HEAD
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
=======
	c.JSON(500, gin.H{
		"status": "OK",
	})
>>>>>>> 83e3dbb9155307f65aef0c2bea5438d96e45d3c5
}
