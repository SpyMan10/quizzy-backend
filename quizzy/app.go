package quizzy

import (
	"github.com/gin-gonic/gin"
	"log"
	"quizzy.app/backend/quizzy/cfg"
	quizzyhttp "quizzy.app/backend/quizzy/http"
	"quizzy.app/backend/quizzy/services"
)

func Run() {
	config := cfg.LoadCfgFromEnv()

	// Configure GIN execution mode (dev, test, production).
	setGinMode(config.Env)

	log.Printf("application running in %s mode.\n", config.Env)

	// Initializing GIN engine.
	engine := gin.Default()
	//engine.Use(cors.Default())

	router := engine.Group(config.BasePath)

	// Configure database provider.
	// Firebase access is injected here into GIN context,
	// this will enable fast access to database through handling chain itself.
	router.Use(func(ctx *gin.Context) {
		//FIXME: Firebase application must be initialized outside ConfigureFirebase().
		// Firestore can be initialized each time we need it.
		if client, err := services.ConfigureFirebase(config); err == nil {
			ctx.Set("firebase-services", client)
		}
	})
	router.Use(func(ctx *gin.Context) {
		redis := services.ConfigureRedis(config)
		if redis != nil {
			ctx.Set("redis-service", redis)
		}
	})

	// Initializing HTTP routes.
	quizzyhttp.ConfigureRouting(router)

	// Running server...
	if err := engine.Run(config.Addr); err != nil {
		log.Fatalf("Failed to start server on %s: %s", config.Addr, err)
	}
}

func setGinMode(env string) {
	switch env {
	case cfg.EnvDevelopment:
		gin.SetMode(gin.DebugMode)
		break
	case cfg.EnvTest:
		gin.SetMode(gin.TestMode)
		break
	default:
		gin.SetMode(gin.ReleaseMode)
		break
	}
}
