package users

import (
	"github.com/gin-gonic/gin"
	svc "quizzy.app/backend/quizzy/services"
)

func useUserStore(ctx *gin.Context) UserStore {
	return ctx.MustGet("user-store").(UserStore)
}

func provideUserStore(ctx *gin.Context) {
	fb := ctx.MustGet("firebase-services").(svc.FirebaseServices)

	if fb.Store != nil {
		ctx.Set("user-store", ConfigureUserStore(fb.Store))
	}
}
