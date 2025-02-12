package cfg

import "os"

const (
	EnvDevelopment = "DEVELOPMENT"
	EnvProduction  = "PRODUCTION"
	EnvTest        = "TEST"
)

type AppConfig struct {
	// Application environment.
	Env string
	// Address to listen on
	Addr string
	// Firebase configuration file.
	FirebaseConfFile string
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
		Env:              getEnvDefault("APP_ENV", EnvProduction),
		Addr:             getEnvDefault("APP_ADDR", ":8000"),
		FirebaseConfFile: os.Getenv("APP_FIREBASE_CONF_FILE"),
	}
}
