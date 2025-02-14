package users

import (
	"github.com/gin-gonic/gin"
	svc "quizzy.app/backend/quizzy/services"
)

func useStore(ctx *gin.Context) Store {
	return ctx.MustGet("user-store").(Store)
}

func provideStore(ctx *gin.Context) {
	fb := ctx.MustGet("firebase-services").(svc.FirebaseServices)

	if fb.Store != nil {
		ctx.Set("user-store", ConfigureStore(fb.Store))
	}
}
