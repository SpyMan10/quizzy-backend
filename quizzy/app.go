package quizzy

import (
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	_ "quizzy.app/backend/docs" // Import Swagger Docs

	"log"
	"quizzy.app/backend/quizzy/auth"
	"quizzy.app/backend/quizzy/cfg"
	"quizzy.app/backend/quizzy/ping"
	"quizzy.app/backend/quizzy/quizzes"
	"quizzy.app/backend/quizzy/services"
	"quizzy.app/backend/quizzy/users"
	"time"
)

// @title Quizzy Backend API
// @version 1.0
// @description API permettant la gestion des utilisateurs et des quiz
// @host localhost:8000
// @BasePath /
// @securityDefinitions.apikey BearerAuth
// @in header
// @name Authorization

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
	engine.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

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
