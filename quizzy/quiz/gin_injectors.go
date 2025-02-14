package quiz

import (
	"errors"
	"github.com/gin-gonic/gin"
	"net/http"
	"quizzy.app/backend/quizzy/middlewares"
	svc "quizzy.app/backend/quizzy/services"
)

func useStore(ctx *gin.Context) Store {
	return ctx.MustGet("quiz-store").(Store)
}

func provideStore(ctx *gin.Context) {
	fb := ctx.MustGet("firebase-services").(svc.FirebaseServices)

	if fb.Store != nil {
		ctx.Set("quiz-store", ConfigureStore(fb.Store))
	}
}

func useQuiz(ctx *gin.Context) Document {
	return ctx.MustGet("current-quiz").(Document)
}

func provideQuiz(ctx *gin.Context) {
	id := middlewares.UseIdentity(ctx)
	store := useStore(ctx)
	qid := ctx.Param("quiz-id")

	if quiz, err := store.GetUnique(id.Uid, qid); err == nil {
		ctx.Set("current-quiz", quiz)
	} else if errors.Is(err, ErrNotFound) {
		ctx.AbortWithStatus(http.StatusNotFound)
	} else {
		ctx.AbortWithStatus(http.StatusInternalServerError)
	}
}
