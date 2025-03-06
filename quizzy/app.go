package quizzy

import (
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
	"log"
	"quizzy.app/backend/quizzy/auth"
	"quizzy.app/backend/quizzy/cfg"
	"quizzy.app/backend/quizzy/ping"
	"quizzy.app/backend/quizzy/quizzes"
	"quizzy.app/backend/quizzy/services"
	"quizzy.app/backend/quizzy/users"
	"time"
)

func Run() {
	config := cfg.LoadCfgFromEnv()

	// Configure GIN execution mode (dev, test, production).
	setGinMode(config.Env.AsString())

	log.Printf("application running in %s mode.\n", config.Env)

	// Initializing GIN engine.
	engine := gin.Default()
	engine.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"http://localhost:4200"},
		AllowMethods:     []string{"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Authorization", "Accept"},
		ExposeHeaders:    []string{"Content-Length", "Location"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}))

	router := engine.Group(config.BasePath)

	fbs, fbsErr := services.ConfigureFirebase(config)
	if fbsErr != nil {
		log.Fatalf("failed to initialize firebase services: %s", fbsErr)
	}

	rc, rcErr := services.ConfigureRedis(config)
	if rcErr != nil {
		log.Fatalf("failed to initialize redis service: %s", rcErr)
	}

	setupModule(router, &fbs, rc, config)

	// Running server...
	if err := engine.Run(config.Addr); err != nil {
		log.Fatalf("Failed to start server on %s: %s\n", config.Addr, err)
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

func setupModule(rt *gin.RouterGroup, fbs *services.FirebaseServices, rc *redis.Client, conf cfg.AppConfig) {
	ping.Configure(fbs, rc).ConfigureRouting(rt)

	secured := rt.Group("", auth.ProvideAuthenticator(&auth.FirebaseAuthenticator{Fbs: fbs}))

	users.Configure(fbs, conf).ConfigureRouting(secured)
	quizzes.Configure(fbs, rc, conf).ConfigureRouting(secured)
}
