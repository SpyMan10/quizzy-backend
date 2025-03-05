package users

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"quizzy.app/backend/quizzy/auth"
)

func ConfigureRoutes(rt *gin.RouterGroup) {
	secured := rt.Group("/users", auth.RequireAuthenticated, ProvideService)
	secured.POST("", handlePostUser)
	secured.GET("/me", handleGetSelf)
}

type CreateUserRequest struct {
	Username string `json:"username"`
}

func handlePostUser(ctx *gin.Context) {
	id := auth.UseIdentity(ctx)

	var req CreateUserRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.AbortWithStatus(http.StatusBadRequest)
		return
	}

	service := UseService(ctx)
	if err := service.Create(User{Username: req.Username, Email: id.Email, Id: id.Uid}); err != nil {
		ctx.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	ctx.Status(http.StatusCreated)
}

func handleGetSelf(ctx *gin.Context) {
	id := auth.UseIdentity(ctx)
	service := UseService(ctx)

	if user, err := service.Get(id.Uid); err == nil {
		ctx.JSON(http.StatusOK, user)
		return
	}

	ctx.AbortWithStatus(http.StatusInternalServerError)
}
