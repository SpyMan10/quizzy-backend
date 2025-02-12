package quiz

import (
	"context"
	"firebase.google.com/go/auth"
	"github.com/gin-gonic/gin"
	"net/http"
	"quizzy.app/backend/quizzy/middlewares"
	"quizzy.app/backend/quizzy/services"
)

//type Quiz struct {
//	Id    string `json:"id"`
//	Title string `json:"title"`
//}

func ConfigureRoutes(router *gin.RouterGroup) {
	secured := router.Group("/quiz", middlewares.RequireAuth)
	secured.GET("", getQuiz)
}

func getQuiz(c *gin.Context) {
	fbs := c.MustGet("firebase-services").(services.FirebaseServices)
	self := c.MustGet("user-token").(auth.Token)

	if doc, err := fbs.Store.Collection("users").Doc(self.UID).Get(context.Background()); err != nil {
		c.AbortWithStatus(http.StatusInternalServerError)
	} else {
		c.JSON(http.StatusOK, gin.H{
			"data": doc.Data()["quiz"],
		})
	}
}
