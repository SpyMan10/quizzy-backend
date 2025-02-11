package main

import (
	"github.com/gin-gonic/gin"
	"log"
	quizzyhttp "quizzy.app/backend/http"
)

func main() {
	// Initializing GIN engine.
	engine := gin.Default()

	// Initializing HTTP routes.
	quizzyhttp.Setup(engine)

	// Running server: listen on any network interface (port 8000).
	if err := engine.Run(":8000"); err != nil {
		log.Fatalf("Error starting server: %s", err)
	}
}
