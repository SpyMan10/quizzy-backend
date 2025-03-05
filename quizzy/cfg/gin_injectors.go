package cfg

import "github.com/gin-gonic/gin"

const KeyAppConfig = "app:global-config"

// UseConfig require in the given handler/middleware.
// If config wasn't defined before, this function will cause a thread panic.
func UseConfig(ctx *gin.Context) AppConfig {
	return ctx.MustGet(KeyAppConfig).(AppConfig)
}

// ProvideConfig define application config to the current middleware chain.
func ProvideConfig(config AppConfig) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		ctx.Set(KeyAppConfig, config)
		ctx.Next()
	}
}
