package users

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"quizzy.app/backend/quizzy/auth"
	"quizzy.app/backend/quizzy/cfg"
	"quizzy.app/backend/quizzy/services"
)

type Controller struct {
	Service UserService
}

func Configure(fbs *services.FirebaseServices, conf cfg.AppConfig) *Controller {
	if !conf.Env.IsTest() {
		return &Controller{Service: &UserServiceImpl{Store: &userFirestore{client: fbs.Store}}}
	} else {
		return &Controller{Service: &UserServiceImpl{Store: _newDummyStore()}}
	}
}

func (uc *Controller) ConfigureRouting(rt *gin.RouterGroup) {
	secured := rt.Group("/users", auth.RequireAuthenticated)
	secured.POST("", uc.handlePostUser)
	secured.GET("/me", uc.handleGetSelf)
}

type createUserRequest struct {
	Username string `json:"username"`
}

func (uc *Controller) handlePostUser(ctx *gin.Context) {
	id := auth.UseIdentity(ctx)

	var req createUserRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.AbortWithStatus(http.StatusBadRequest)
		return
	}

	if err := uc.Service.Create(User{Username: req.Username, Email: id.Email, Id: id.Uid}); err != nil {
		ctx.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	ctx.Status(http.StatusCreated)
}

func (uc *Controller) handleGetSelf(ctx *gin.Context) {
	id := auth.UseIdentity(ctx)

	if user, err := uc.Service.Get(id.Uid); err == nil {
		ctx.JSON(http.StatusOK, user)
		return
	}

	ctx.AbortWithStatus(http.StatusInternalServerError)
}
