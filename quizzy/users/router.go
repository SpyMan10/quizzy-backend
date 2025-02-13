package users

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"quizzy.app/backend/quizzy/middlewares"
)

func ConfigureRoutes(rt *gin.RouterGroup) {
	secured := rt.Group("/users", middlewares.RequireAuth, provideUserStore)
	secured.POST("", handlePostUser)
	secured.GET("/me", handleGetSelf)
}

type userRegistrationRequest struct {
	Username string `json:"username"`
}

func handlePostUser(ctx *gin.Context) {
	id := middlewares.UseIdentity(ctx)
	ufc := userRegistrationRequest{}

	if err := ctx.ShouldBindJSON(&ufc); err != nil {
		ctx.AbortWithStatus(http.StatusBadRequest)
		return
	}

	store := useUserStore(ctx)
	if err := store.Upsert(User{Username: ufc.Username, Uid: id.Uid}); err != nil {
		ctx.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	ctx.Status(http.StatusCreated)
}

type UserResponse struct {
	Uid      string `json:"uid"`
	Email    string `json:"email"`
	Username string `json:"username"`
}

func handleGetSelf(ctx *gin.Context) {
	id := middlewares.UseIdentity(ctx)
	store := useUserStore(ctx)

	if user, err := store.GetUnique(id.Uid); err == nil {
		ctx.JSON(http.StatusOK, UserResponse{
			Uid:      user.Uid,
			Email:    id.Email,
			Username: user.Username,
		})
	} else {
		ctx.AbortWithStatus(http.StatusInternalServerError)
	}
}
