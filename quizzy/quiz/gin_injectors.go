package quiz

import (
	"errors"
	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
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
func provideCodeResolver(ctx *gin.Context) {
	rs := ctx.MustGet("redis-service").(*redis.Client)

	if rs != nil {
		ctx.Set("quiz-code-resolver", &redisAdapter{
			client: rs,
		})
	}
}

func useCodeResolver(ctx *gin.Context) QuizCodeResolver {
	return ctx.MustGet("quiz-code-resolver").(QuizCodeResolver)
}

func useQuiz(ctx *gin.Context) Quiz {
	return ctx.MustGet("current-quiz").(Quiz)
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

func provideQuestion(ctx *gin.Context) {
	id := middlewares.UseIdentity(ctx)
	store := useStore(ctx)
	qid := ctx.Param("question-id")
	quiz := useQuiz(ctx)

	if question, err := store.GetUniqueQuestion(id.Uid, quiz.Id, qid); err == nil {
		ctx.Set("current-question", question)
	} else if errors.Is(err, ErrNotFound) {
		ctx.AbortWithStatus(http.StatusNotFound)
	} else {
		ctx.AbortWithStatus(http.StatusInternalServerError)
	}
}

func useQuestion(ctx *gin.Context) Question {
	return ctx.MustGet("current-question").(Question)
}
