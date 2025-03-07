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
	return &Controller{Service: &UserServiceImpl{Store: &userFirestore{client: fbs.Store}}}
}

func (uc *Controller) ConfigureRouting(rt *gin.RouterGroup) {
	secured := rt.Group("/users", auth.RequireAuthenticated)
	secured.POST("", uc.handlePostUser)
	secured.GET("/me", uc.handleGetSelf)
}

type createUserRequest struct {
	Username string `json:"username"`
}

// handlePostUser crée un nouvel utilisateur
// @Summary Créer un utilisateur
// @Description Cette route permet de créer un nouvel utilisateur à partir d'un username et de l'email récupéré via l'authentification
// @Tags Users
// @Accept json
// @Produce json
// @Param Authorization header string true "Token d'authentification Bearer" default(Bearer <votre_token>)
// @Param body body createUserRequest true "Informations de l'utilisateur à créer"
// @Success 201 {string} string "Utilisateur créé avec succès"
// @Failure 400 {string} string "Requête invalide"
// @Failure 500 {string} string "Erreur interne du serveur"
// @Router /users [post]
// @Security BearerAuth
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

// handleGetSelf récupère les informations de l'utilisateur authentifié
// @Summary Récupérer les informations de l'utilisateur connecté
// @Description Cette route permet d'obtenir les informations du compte actuellement authentifié
// @Tags Users
// @Produce json
// @Param Authorization header string true "Token d'authentification Bearer" default(Bearer <votre_token>)
// @Success 200 {object} User "Informations de l'utilisateur"
// @Failure 401 {string} string "Utilisateur non authentifié"
// @Failure 500 {string} string "Erreur interne du serveur"
// @Router /users/me [get]
// @Security BearerAuth
func (uc *Controller) handleGetSelf(ctx *gin.Context) {
	id := auth.UseIdentity(ctx)

	if user, err := uc.Service.Get(id.Uid); err == nil {
		ctx.JSON(http.StatusOK, user)
		return
	}

	ctx.AbortWithStatus(http.StatusInternalServerError)
}
