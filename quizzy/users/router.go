package users

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"quizzy.app/backend/quizzy/middlewares"
)

func ConfigureRoutes(rt *gin.RouterGroup) {
	secured := rt.Group("/users", middlewares.RequireAuth, provideStore)
	secured.POST("", handlePostUser)
	secured.GET("/me", handleGetSelf)
}

type CreateUserRequest struct {
	Username string `json:"username"`
}

func handlePostUser(ctx *gin.Context) {
	id := middlewares.UseIdentity(ctx)

	var req CreateUserRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.AbortWithStatus(http.StatusBadRequest)
		return
	}

	store := useStore(ctx)
	if err := store.Upsert(Document{Username: req.Username, Email: id.Email, Uid: id.Uid}); err != nil {
		ctx.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	ctx.Status(http.StatusCreated)
}

func handleGetSelf(ctx *gin.Context) {
	id := middlewares.UseIdentity(ctx)
	store := useStore(ctx)

	if user, err := store.GetUnique(id.Uid); err == nil {
		ctx.JSON(http.StatusOK, user)
		return
	}

	ctx.AbortWithStatus(http.StatusInternalServerError)
}
