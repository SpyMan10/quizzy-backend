package quiz

import (
	"context"
	"firebase.google.com/go/auth"
	"github.com/gin-gonic/gin"
	"net/http"
	"quizzy.app/backend/quizzy/middlewares"
	"quizzy.app/backend/quizzy/services"
)

func ConfigureRoutes(router *gin.RouterGroup) {
	secured := router.Group("/quizz", middlewares.RequireAuth)
	secured.GET("", getAllUserQuiz)
}

func getAllUserQuiz(c *gin.Context) {
	fbs := c.MustGet("firebase-services").(services.FirebaseServices)
	self := c.MustGet("user-token").(*auth.Token)

	if doc, err := fbs.Store.Collection("users").Doc(self.UID).Get(context.Background()); err != nil {
		c.AbortWithStatus(http.StatusInternalServerError)
	} else {
		if doc.Data()["quiz"] == nil {
			c.JSON(http.StatusOK, gin.H{
				"data": []any{},
			})
			return
		}
		c.JSON(http.StatusOK, gin.H{
			"data": doc.Data()["quiz"],
		})
	}
}
