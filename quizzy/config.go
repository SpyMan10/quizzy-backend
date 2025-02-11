package quizzy

import "os"

const (
	EnvDevelopment = "DEVELOPMENT"
	EnvProduction  = "PRODUCTION"
)

type AppConfig struct {
	// Application environment.
	env string
	// Address to listen on
	addr string
	// Firebase configuration file.
	firebaseConfFile string
}

// getEnvDefault fetch environment variable from the given key and return it if found,
// otherwise the default value is returned.
func getEnvDefault(key, def string) string {
	if v, f := os.LookupEnv(key); f {
		return v
	}

	return def
}

func LoadCfgFromEnv() AppConfig {
	return AppConfig{
		env:              getEnvDefault("APP_ENV", EnvProduction),
		addr:             getEnvDefault("APP_ADDR", ":8000"),
		firebaseConfFile: os.Getenv("APP_FIREBASE_CONF_FILE"),
	}
}
