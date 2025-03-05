package users

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"quizzy.app/backend/quizzy/cfg"
	svc "quizzy.app/backend/quizzy/services"
)

const KeyUserService = "user-service"

func UseService(ctx *gin.Context) UserService {
	return ctx.MustGet(KeyUserService).(UserService)
}

func ProvideService(ctx *gin.Context) {
	conf := cfg.UseConfig(ctx)

	if !conf.Env.IsTest() {
		fb := ctx.MustGet("firebase-services").(svc.FirebaseServices)

		if fb.Store != nil {
			ctx.Set(KeyUserService, &UserServiceImpl{Store: ConfigureStore(fb.Store)})
		}

		fmt.Printf("info: user-service plugged in to firestore.\n")
	} else {
		ctx.Set(KeyUserService, DummyUserStoreImpl{})
		fmt.Printf("info: user-service plugged in to dummy impl.\n")
	}
}
