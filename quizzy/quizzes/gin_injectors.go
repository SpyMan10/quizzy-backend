package quizzes

import (
	"errors"
	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
	"net/http"
	"quizzy.app/backend/quizzy/auth"
	"quizzy.app/backend/quizzy/cfg"
	svc "quizzy.app/backend/quizzy/services"
)

const KeyQuizService = "quiz-service"

func UseService(ctx *gin.Context) QuizService {
	return ctx.MustGet(KeyQuizService).(QuizService)
}

func ProvideService(ctx *gin.Context) {
	fb := ctx.MustGet("firebase-services").(svc.FirebaseServices)
	rc := ctx.MustGet("redis-service").(*redis.Client)
	conf := cfg.UseConfig(ctx)

	if conf.Env.IsTest() {
		ctx.Set(KeyQuizService, &QuizServiceImpl{
			store:    &DummyQuizStoreImpl{entries: make([]dummyEntry, 0)},
			resolver: &dummyCodeResolver{entries: make(map[string]string)},
		})
	} else {
		if fb.Store != nil {
			ctx.Set(KeyQuizService, &QuizServiceImpl{
				store:    &fireStoreAdapter{client: fb.Store},
				resolver: &RedisCodeResolver{client: rc},
			})
		}
	}
}

func UseCodeResolver(ctx *gin.Context) QuizCodeResolver {
	return ctx.MustGet("quiz-code-resolver").(QuizCodeResolver)
}

func UseQuiz(ctx *gin.Context) Quiz {
	return ctx.MustGet("current-quiz").(Quiz)
}

func ProvideQuiz(ctx *gin.Context) {
	id := auth.UseIdentity(ctx)
	service := UseService(ctx)
	qid := ctx.Param("quiz-id")

	if quiz, err := service.Get(id.Uid, qid); err == nil {
		ctx.Set("current-quiz", quiz)
	} else if errors.Is(err, ErrNotFound) {
		ctx.AbortWithStatus(http.StatusNotFound)
	} else {
		ctx.AbortWithStatus(http.StatusInternalServerError)
	}
}

func ProvideQuestion(ctx *gin.Context) {
	qid := ctx.Param("question-id")
	quiz := UseQuiz(ctx)

	for _, q := range quiz.Questions {
		if q.Id == qid {
			ctx.Set("current-question", q)
			return
		}
	}

	ctx.AbortWithStatus(http.StatusNotFound)
}

func UseQuestion(ctx *gin.Context) Question {
	return ctx.MustGet("current-question").(Question)
}
