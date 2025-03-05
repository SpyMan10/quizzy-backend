package auth

import (
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
	"strings"
)

const (
	KeyAuthenticator = "authenticator"
	KeyIdentity      = "identity"
)

// RequireAuthenticated middleware perform authorization against for the current request.
// If the given authorization isn't valid, this middleware will stop propagation, and immediately
// abort request processing.
// If authorization succeed, a new Identity will be injected in the current middleware chain.
func RequireAuthenticated(ctx *gin.Context) {
	token := strings.TrimSpace(strings.TrimLeft(ctx.GetHeader("Authorization"), "Bearer"))

	if len(token) == 0 {
		log.Println("missing authorization token")
		ctx.AbortWithStatus(http.StatusUnauthorized)
		return
	}

	authenticator := UseAuthenticator(ctx)
	if id, err := authenticator.Authorize(token); err != nil {
		ctx.AbortWithStatus(http.StatusUnauthorized)
	} else {
		ctx.Set(KeyIdentity, id)
		ctx.Next()
	}
}

func UseIdentity(ctx *gin.Context) Identity {
	return ctx.MustGet(KeyIdentity).(Identity)
}

// ProvideAuthenticator expose the given authenticator to the current middleware chain.
func ProvideAuthenticator(authenticator Authenticator) gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Set(KeyAuthenticator, authenticator)
	}
}

func UseAuthenticator(ctx *gin.Context) Authenticator {
	return ctx.MustGet(KeyAuthenticator).(Authenticator)
}
