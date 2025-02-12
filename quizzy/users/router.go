package users

import (
	"context"
	"firebase.google.com/go/auth"
	"github.com/gin-gonic/gin"
	svc "quizzy.app/backend/quizzy/services"
	"strings"
)

func RequireAuth(ctx *gin.Context) {
	token := strings.TrimLeft(ctx.GetHeader("Authorization"), "Bearer ")
	services := ctx.MustGet("firebase-services").(svc.FirebaseServices)

	if len(token) == 0 {
		ctx.AbortWithStatus(401)
		return
	}

	if tok, err := services.Auth.VerifyIDTokenAndCheckRevoked(context.Background(), token); err != nil {
		ctx.Set("user-token", tok)
		ctx.Next()
	}
}

func ConfigureRoutes(rt *gin.RouterGroup) {
	rt.POST("/users", RequireAuth, createUser)
}

type userForCreate struct {
	Username string `json:"username"`
}

func createUser(c *gin.Context) {
	servicesFire := c.MustGet("firebase-services").(svc.FirebaseServices)
	userToken := c.MustGet("user-token").(*auth.Token)

	ufc := userForCreate{}
	if err := c.ShouldBindJSON(&ufc); err != nil {
		c.AbortWithStatus(400)
		return
	}
	_, err := servicesFire.Store.Collection("users").Doc(userToken.UID).Set(context.Background(), map[string]interface{}{
		"username": ufc.Username,
	})
	if err != nil {
		return
	}
}
