package users

import (
	"context"
	"firebase.google.com/go/auth"
	"github.com/gin-gonic/gin"
	"quizzy.app/backend/quizzy/middlewares"
	svc "quizzy.app/backend/quizzy/services"
)

func ConfigureRoutes(rt *gin.RouterGroup) {
	secured := rt.Group("/users", middlewares.RequireAuth)
	secured.POST("", postUser)
	secured.GET("/me", getSelf)
}

type userForCreate struct {
	Username string `json:"username"`
}

func postUser(c *gin.Context) {
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
		c.AbortWithStatus(200)
	}

	c.Status(200)
}

type UserData struct {
	Username string `json:"username"`
	Email    string `json:"email"`
	Uid      string `json:"uid"`
}

func getSelf(c *gin.Context) {
	servicesFire := c.MustGet("firebase-services").(svc.FirebaseServices)
	userToken := c.MustGet("user-token").(*auth.Token)

	if doc, err := servicesFire.Store.Collection("users").Doc(userToken.UID).Get(context.Background()); err != nil {
		c.AbortWithStatus(500)
	} else {
		c.JSON(200, UserData{
			Username: doc.Data()["username"].(string),
			Uid:      userToken.UID,
			Email:    userToken.Claims["email"].(string),
		})
	}
}
