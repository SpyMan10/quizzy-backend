package middlewares

import (
	"context"
	"github.com/gin-gonic/gin"
	"log"
	svc "quizzy.app/backend/quizzy/services"
	"strings"
)

type Identity struct {
	Uid   string
	Email string
}

func RequireAuth(ctx *gin.Context) {
	token := strings.TrimSpace(strings.TrimLeft(ctx.GetHeader("Authorization"), "Bearer"))

	if len(token) == 0 {
		log.Println("missing authorization token")
		ctx.AbortWithStatus(401)
		return
	}

	services := ctx.MustGet("firebase-services").(svc.FirebaseServices)
	if tok, err := services.Auth.VerifyIDTokenAndCheckRevoked(context.Background(), token); err != nil {
		ctx.AbortWithStatus(401)
	} else {
		ctx.Set("user-token", tok)
		ctx.Set("identity", Identity{
			Uid:   tok.UID,
			Email: tok.Claims["email"].(string),
		})
		ctx.Next()
	}
}

func UseIdentity(ctx *gin.Context) Identity {
	return ctx.MustGet("identity").(Identity)
}
