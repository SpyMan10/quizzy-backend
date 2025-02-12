package quizzy

import (
	"github.com/gin-gonic/gin"
	"log"
	cfg2 "quizzy.app/backend/quizzy/cfg"
	quizzyhttp "quizzy.app/backend/quizzy/http"
)

func Run() {
	cfg := cfg2.LoadCfgFromEnv()

	switch cfg.env {
	case cfg2.EnvDevelopment:
		gin.SetMode(gin.DebugMode)
		break
	case cfg2.EnvTest:
		gin.SetMode(gin.TestMode)
		break
	default:
		gin.SetMode(gin.ReleaseMode)
		break
	}

	log.Printf("running mode : %s\n", cfg.env)

	// Initializing GIN engine.
	engine := gin.Default()

	// Creating base router.
	router := engine.Group("/")

	// Configure database provider.
	// Firebase access is injected here into GIN context,
	// this will enable fast access to database through handling chain itself.
	router.Use(func(ctx *gin.Context) {
		//FIXME: Firebase application must be initialized outside ConfigureFirebase().
		// Firestore can be initialized each time we need it.
		if client, err := ConfigureFirebase(cfg); err == nil && client != nil {
			println("Firebase initialized.")
			ctx.Set("firebase-services", client)
		}
	})

	// Initializing HTTP routes.
	quizzyhttp.ConfigureRouting(router)

	// Running server...
	if err := engine.Run(cfg.addr); err != nil {
		log.Fatalf("Failed to start server on %s: %s", cfg.addr, err)
	}
}
