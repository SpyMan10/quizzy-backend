package quizzes

import (
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/redis/go-redis/v9"
	"net/http"
	"quizzy.app/backend/quizzy/auth"
	"quizzy.app/backend/quizzy/cfg"
	"quizzy.app/backend/quizzy/services"
)

type Controller struct {
	Resolver QuizCodeResolver
	Service  QuizService
}

func Configure(fbs *services.FirebaseServices, rc *redis.Client, conf cfg.AppConfig) *Controller {
	if !conf.Env.IsTest() {
		return &Controller{
			Service: &QuizServiceImpl{
				store:    &quizFirestore{client: fbs.Store},
				resolver: &RedisCodeResolver{client: rc},
			},
		}
	} else {
		return &Controller{
			Service: &QuizServiceImpl{
				store:    _createDummyStore(),
				resolver: &dummyCodeResolver{entries: make(map[string]string)},
			},
		}
	}
}

func (qc *Controller) ConfigureRouting(rt *gin.RouterGroup) {
	secured := rt.Group("/quiz", auth.RequireAuthenticated)
	secured.GET("", qc.handleGetAllUserQuiz)
	secured.POST("", qc.handlePostQuiz)

	quiz := secured.Group("/:quiz-id", qc.ProvideQuiz)
	quiz.GET("", handleGetQuiz)
	quiz.PATCH("", qc.handlePatchQuiz)
	quiz.GET("/questions", handleGetQuestions)
	quiz.POST("/questions", qc.handlePostQuestion)

	quiz.PUT("/questions/:question-id", ProvideQuestion, qc.handlePutQuestion)
	quiz.POST("/start", qc.handleStartQuiz)
}

func UseQuiz(ctx *gin.Context) Quiz {
	return ctx.MustGet("current-quiz").(Quiz)
}

func (qc *Controller) ProvideQuiz(ctx *gin.Context) {
	id := auth.UseIdentity(ctx)
	qid := ctx.Param("quiz-id")

	if quiz, err := qc.Service.Get(id.Uid, qid); err == nil {
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

// handleGetQuiz retourne les détails d'un quiz spécifique
// @Summary Récupérer un quiz
// @Description Retourne les informations d'un quiz appartenant à l'utilisateur authentifié
// @Tags Quizzes
// @Produce json
// @Param Authorization header string true "Token d'authentification Bearer" default(Bearer <votre_token>)
// @Param quiz-id path string true "ID du quiz"
// @Success 200 {object} Quiz "Détails du quiz"
// @Failure 401 {string} string "Utilisateur non authentifié"
// @Failure 404 {string} string "Quiz non trouvé"
// @Failure 500 {string} string "Erreur interne du serveur"
// @Router /quiz/{quiz-id} [get]
// @Security BearerAuth
func handleGetQuiz(ctx *gin.Context) {
	quiz := UseQuiz(ctx)
	ctx.JSON(http.StatusOK, quiz)
}

type QuizWithLinks struct {
	Quiz
	Links Links `json:"_links"`
}

type UserQuizzesResponse struct {
	Data  []QuizWithLinks `json:"data"`
	Links Links           `json:"_links"`
}

func mapMultipleQuizWithLinks(quizzes []Quiz) []QuizWithLinks {
	qwl := make([]QuizWithLinks, len(quizzes))

	for i := range quizzes {
		qwl[i] = mapQuizWithLinks(quizzes[i])
	}

	return qwl
}

func mapQuizWithLinks(quiz Quiz) QuizWithLinks {
	lnk := Links{
		Create: "",
		Start:  "",
	}

	if quiz.Validate() {
		lnk.Start = fmt.Sprintf("http://localhost:8000/quiz/%s/start", quiz.Id)
	}

	return QuizWithLinks{
		Quiz:  quiz,
		Links: lnk,
	}
}

// handleGetAllUserQuiz retourne tous les quiz de l'utilisateur connecté
// @Summary Récupérer tous mes quiz
// @Description Retourne la liste des quiz créés par l'utilisateur authentifié
// @Tags Quizzes
// @Produce json
// @Param Authorization header string true "Token d'authentification Bearer" default(Bearer <votre_token>)
// @Success 200 {object} UserQuizzesResponse "Liste des quiz de l'utilisateur"
// @Failure 401 {string} string "Utilisateur non authentifié"
// @Failure 500 {string} string "Erreur interne du serveur"
// @Router /quiz [get]
// @Security BearerAuth
func (qc *Controller) handleGetAllUserQuiz(ctx *gin.Context) {
	id := auth.UseIdentity(ctx)

	if quizzes, err := qc.Service.GetAll(id.Uid); err == nil {
		ctx.JSON(http.StatusOK, UserQuizzesResponse{
			Data: mapMultipleQuizWithLinks(quizzes),
			Links: Links{
				Create: "http://localhost:8000/quiz",
			},
		})
		return
	}

	ctx.AbortWithStatus(http.StatusInternalServerError)
}

type CreateQuizRequest struct {
	Title       string `json:"title"`
	Description string `json:"description"`
}

// handlePostQuiz crée un nouveau quiz
// @Summary Créer un quiz
// @Description Permet à l'utilisateur authentifié de créer un quiz
// @Tags Quizzes
// @Accept json
// @Produce json
// @Param Authorization header string true "Token d'authentification Bearer" default(Bearer <votre_token>)
// @Param body body CreateQuizRequest true "Informations du quiz à créer"
// @Success 201 {object} Quiz "Quiz créé avec succès"
// @Failure 400 {string} string "Requête invalide"
// @Failure 401 {string} string "Utilisateur non authentifié"
// @Failure 500 {string} string "Erreur interne du serveur"
// @Router /quiz [post]
// @Security BearerAuth
func (qc *Controller) handlePostQuiz(ctx *gin.Context) {
	id := auth.UseIdentity(ctx)

	// Getting payload from request.
	var req CreateQuizRequest
	if ctx.ShouldBindJSON(&req) != nil {
		ctx.AbortWithStatus(http.StatusBadRequest)
		return
	}

	if code, err := GenerateCode(); err == nil {
		quiz := Quiz{
			Id:          uuid.New().String(),
			Title:       req.Title,
			Description: req.Description,
			Code:        code,
		}

		if err2 := qc.Service.Create(id.Uid, quiz); err2 == nil {
			ctx.Header("Location", fmt.Sprintf("http://localhost:8000/quiz/%s", quiz.Id))
			ctx.JSON(http.StatusCreated, quiz)
			return
		}
	}

	// WARNING / WARNING / WARNING //
	// If this code-path is reached, it means that the requested user was never
	// registered in our firestore.
	// WARNING / WARNING / WARNING //

	// This may happen if client does some weird things...
	// We should never let the client decide about this process,
	// user registration must be done in single request (Client->Server), or we must use pub/sub (or something similar)
	// from firebase to service user automatically to avoid data consistency issues.
	ctx.AbortWithStatus(http.StatusInternalServerError)
}

type PatchQuizRequest []FieldPatchOp

// handlePatchQuiz met à jour un quiz existant
// @Summary Modifier un quiz
// @Description Met à jour un quiz existant en fonction des champs envoyés
// @Tags Quizzes
// @Accept json
// @Produce json
// @Param Authorization header string true "Token d'authentification Bearer" default(Bearer <votre_token>)
// @Param quiz-id path string true "ID du quiz"
// @Param body body PatchQuizRequest true "Champs à modifier"
// @Success 204 {string} string "Quiz mis à jour avec succès"
// @Failure 400 {string} string "Requête invalide"
// @Failure 401 {string} string "Utilisateur non authentifié"
// @Failure 500 {string} string "Erreur interne du serveur"
// @Router /quiz/{quiz-id} [patch]
// @Security BearerAuth
func (qc *Controller) handlePatchQuiz(ctx *gin.Context) {
	id := auth.UseIdentity(ctx)
	quiz := UseQuiz(ctx)

	var req PatchQuizRequest
	if ctx.ShouldBindJSON(&req) != nil {
		ctx.AbortWithStatus(http.StatusBadRequest)
		return
	}

	if err := qc.Service.Patch(id.Uid, quiz.Id, req); err == nil {
		ctx.Status(http.StatusNoContent)
	} else if errors.Is(err, ErrInvalidPatchOperator) || errors.Is(err, ErrInvalidPatchField) {
		ctx.AbortWithStatus(http.StatusBadRequest)
	} else {
		ctx.AbortWithStatus(http.StatusInternalServerError)
	}
}

type CreateQuestionRequest struct {
	Title   string   `json:"title"`
	Answers []Answer `json:"answers"`
}

// handlePostQuestion ajoute une question à un quiz
// @Summary Ajouter une question
// @Description Permet d'ajouter une question à un quiz existant
// @Tags Quizzes
// @Accept json
// @Produce json
// @Param Authorization header string true "Token d'authentification Bearer" default(Bearer <votre_token>)
// @Param quiz-id path string true "ID du quiz"
// @Param body body CreateQuestionRequest true "Détails de la question"
// @Success 201 {string} string "Question ajoutée avec succès"
// @Failure 400 {string} string "Requête invalide"
// @Failure 401 {string} string "Utilisateur non authentifié"
// @Failure 500 {string} string "Erreur interne du serveur"
// @Router /quiz/{quiz-id}/questions [post]
// @Security BearerAuth
func (qc *Controller) handlePostQuestion(ctx *gin.Context) {
	id := auth.UseIdentity(ctx)
	quiz := UseQuiz(ctx)

	var req CreateQuestionRequest
	if ctx.ShouldBindJSON(&req) != nil {
		ctx.AbortWithStatus(http.StatusBadRequest)
		return
	}

	question := Question{
		Id:      uuid.New().String(),
		Title:   req.Title,
		Answers: req.Answers,
	}
	err := qc.Service.CreateQuestion(id.Uid, quiz, question)

	if err != nil {
		ctx.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	ctx.Header("Location", fmt.Sprintf("http://localhost:8000/quiz/%s/questions/%s", quiz.Id, question.Id))
	ctx.Status(http.StatusCreated)
}

// handleGetQuestions retourne la liste des questions d'un quiz
// @Summary Récupérer les questions d'un quiz
// @Description Retourne toutes les questions du quiz spécifié par son ID
// @Tags Quizzes
// @Produce json
// @Param Authorization header string true "Token d'authentification Bearer" default(Bearer <votre_token>)
// @Param quiz-id path string true "ID du quiz"
// @Success 200 {array} Question "Liste des questions du quiz"
// @Failure 401 {string} string "Utilisateur non authentifié"
// @Failure 404 {string} string "Quiz non trouvé"
// @Failure 500 {string} string "Erreur interne du serveur"
// @Router /quiz/{quiz-id}/questions [get]
// @Security BearerAuth
func handleGetQuestions(ctx *gin.Context) {
	quiz := UseQuiz(ctx)
	ctx.JSON(http.StatusOK, quiz.Questions)
}

type UnidentifiedAnswer struct {
	Title     string `json:"title"`
	IsCorrect bool   `json:"isCorrect"`
}

type UpdateQuestionRequest struct {
	Title   string               `json:"title"`
	Answers []UnidentifiedAnswer `json:"answers"`
}

// handlePutQuestion met à jour une question existante
// @Summary Modifier une question
// @Description Met à jour une question spécifique d'un quiz
// @Tags Quizzes
// @Accept json
// @Produce json
// @Param Authorization header string true "Token d'authentification Bearer" default(Bearer <votre_token>)
// @Param quiz-id path string true "ID du quiz"
// @Param question-id path string true "ID de la question"
// @Param body body UpdateQuestionRequest true "Mise à jour de la question"
// @Success 204 {string} string "Question mise à jour avec succès"
// @Failure 400 {string} string "Requête invalide"
// @Failure 401 {string} string "Utilisateur non authentifié"
// @Failure 404 {string} string "Quiz ou question non trouvée"
// @Failure 500 {string} string "Erreur interne du serveur"
// @Router /quiz/{quiz-id}/questions/{question-id} [put]
// @Security BearerAuth
func (qc *Controller) handlePutQuestion(ctx *gin.Context) {
	id := auth.UseIdentity(ctx)
	quiz := UseQuiz(ctx)
	question := UseQuestion(ctx)

	var payload UpdateQuestionRequest
	if err := ctx.ShouldBindJSON(&payload); err != nil {
		ctx.AbortWithStatus(http.StatusBadRequest)
		return
	}

	question.Title = payload.Title
	question.Answers = make([]Answer, 0)
	for _, a := range payload.Answers {
		question.Answers = append(question.Answers, Answer{
			Id:        uuid.New().String(),
			Title:     a.Title,
			IsCorrect: a.IsCorrect,
		})
	}

	if err := qc.Service.UpdateQuestion(id.Uid, quiz.Id, question); err != nil {
		ctx.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	ctx.Status(http.StatusNoContent)
}

// handleStartQuiz démarre un quiz
// @Summary Démarrer un quiz
// @Description Démarre un quiz et retourne son code d'exécution
// @Tags Quizzes
// @Produce json
// @Param Authorization header string true "Token d'authentification Bearer" default(Bearer <votre_token>)
// @Param quiz-id path string true "ID du quiz"
// @Success 201 {string} string "Quiz démarré avec succès"
// @Failure 400 {string} string "Quiz non prêt à être démarré"
// @Failure 401 {string} string "Utilisateur non authentifié"
// @Failure 404 {string} string "Quiz non trouvé"
// @Failure 500 {string} string "Erreur interne du serveur"
// @Router /quiz/{quiz-id}/start [post]
// @Security BearerAuth
func (qc *Controller) handleStartQuiz(ctx *gin.Context) {
	identity := auth.UseIdentity(ctx)
	quiz := UseQuiz(ctx)

	if err := qc.Service.StartQuiz(identity.Uid, quiz); errors.Is(err, ErrQuizNotReady) {
		ctx.AbortWithStatus(http.StatusBadRequest)
		return
	} else if err != nil {
		ctx.AbortWithStatus(http.StatusBadRequest)
		return
	}

	ctx.Header("Location", fmt.Sprintf("http://localhost:8000/execution/%s", quiz.Code))
	ctx.Status(http.StatusCreated)
}
