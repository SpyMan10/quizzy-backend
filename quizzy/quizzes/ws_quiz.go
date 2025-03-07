package quizzes

import (
	"encoding/json"
	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"log"
	"net/http"
	"sync"
)

type SocketController struct {
	Service       QuizService
	upgrade       websocket.Upgrader
	rooms         map[string][]*websocket.Conn
	hosts         map[string]*websocket.Conn
	roomsMu       sync.Mutex
	questionIdx   map[string]int
	questionIdxMu sync.Mutex
}

func NewSocketController(service QuizService) *SocketController {
	return &SocketController{
		Service: service,
		upgrade: websocket.Upgrader{
			ReadBufferSize:  1024,
			WriteBufferSize: 1024,
			CheckOrigin: func(r *http.Request) bool {
				return true
			},
		},
		rooms:       make(map[string][]*websocket.Conn),
		hosts:       make(map[string]*websocket.Conn),
		questionIdx: make(map[string]int),
	}
}

// configureWs initialise la connexion WebSocket
// @Summary Connexion WebSocket
// @Description Établit une connexion WebSocket pour interagir avec le quiz en temps réel
// @Tags WebSocket
// @Produce json
// @Param Authorization header string true "Token d'authentification Bearer" default(Bearer <votre_token>)
// @Success 101 {string} string "Connexion WebSocket établie"
// @Failure 400 {string} string "Mauvaise requête"
// @Failure 401 {string} string "Non authentifié"
// @Router /quiz/ws [get]
// @Security BearerAuth
func (sc *SocketController) Configure(router *gin.RouterGroup) {
	router.GET("/", func(c *gin.Context) {
		sc.handleWebSocket(c.Writer, c.Request)
	})
}

func (sc *SocketController) handleWebSocket(w http.ResponseWriter, r *http.Request) {
	conn, err := sc.upgrade.Upgrade(w, r, nil)
	if err != nil {
		log.Println("Failed to set websocket upgrade: ", err)
		return
	}

	for {
		_, msg, err := conn.ReadMessage()
		if err != nil {
			log.Println("Error reading message:", err)
			break
		}

		var event map[string]interface{}
		if err := json.Unmarshal(msg, &event); err != nil {
			log.Println("Error unmarshalling message:", err)
			break
		}

		switch event["name"] {
		case "host":
			sc.handleHostEvent(conn, event["data"].(map[string]any))
		case "join":
			sc.handleJoinEvent(conn, event["data"].(map[string]any))
		case "nextQuestion":
			sc.handleNextQuestionEvent(event["data"].(map[string]any))
		}
	}
}

// handleHostEvent permet d'héberger un quiz via WebSocket
// @Summary Héberger un quiz
// @Description L'utilisateur devient l'hôte d'un quiz et reçoit les détails du quiz
// @Tags WebSocket
// @Accept json
// @Produce json
// @Param event body object true "Événement WebSocket 'host'"
// @Success 200 {object} map[string]interface{} "Détails du quiz et du statut"
// @Router /quiz/ws [post]
// @Security BearerAuth

func (sc *SocketController) handleHostEvent(conn *websocket.Conn, data map[string]any) {
	executionId := data["executionId"].(string)
	sc.roomsMu.Lock()
	sc.hosts[executionId] = conn                // On stocke l'host séparément
	sc.rooms[executionId] = []*websocket.Conn{} // On initialise la room sans participants
	sc.roomsMu.Unlock()

	quiz, err := sc.Service.QuizFromCode(executionId)
	if err != nil {
		return
	}

	_ = sc.Service.ResetRoomPeople(executionId)

	response := map[string]interface{}{
		"name": "hostDetails",
		"data": map[string]interface{}{
			"quiz": quiz.Title,
		},
	}
	res, _ := json.Marshal(response)
	_ = conn.WriteMessage(websocket.TextMessage, res)

	// Exclure l'hôte du comptage
	nbPeoples, _ := sc.Service.GetRoomPeople(executionId)

	sc.broadcastToRoom(executionId, map[string]interface{}{
		"name": "status",
		"data": map[string]interface{}{
			"status":       "waiting",
			"participants": nbPeoples, // Pas d'incrémentation pour l'host
		},
	})
	sc.questionIdxMu.Lock()
	sc.questionIdx[executionId] = 0 // Initialiser l'index des questions
	sc.questionIdxMu.Unlock()
}

// handleJoinEvent permet à un utilisateur de rejoindre un quiz via WebSocket
// @Summary Rejoindre un quiz
// @Description Un utilisateur rejoint un quiz et reçoit les détails du quiz
// @Tags WebSocket
// @Accept json
// @Produce json
// @Param event body object true "Événement WebSocket 'join'"
// @Success 200 {object} map[string]interface{} "Détails du quiz et du statut"
// @Router /quiz/ws [post]
// @Security BearerAuth

func (sc *SocketController) handleJoinEvent(conn *websocket.Conn, data map[string]any) {
	executionId := data["executionId"].(string)
	sc.roomsMu.Lock()
	sc.rooms[executionId] = append(sc.rooms[executionId], conn) // Ajout uniquement aux participants
	sc.roomsMu.Unlock()

	quiz, err := sc.Service.QuizFromCode(executionId)
	if err != nil {
		return
	}
	_ = sc.Service.IncrRoomPeople(executionId) // On incrémente uniquement pour les participants
	nbPeoples, _ := sc.Service.GetRoomPeople(executionId)

	response := map[string]interface{}{
		"name": "joinDetails",
		"data": map[string]interface{}{
			"quizTitle": quiz.Title,
		},
	}
	res, _ := json.Marshal(response)
	_ = conn.WriteMessage(websocket.TextMessage, res)

	sc.broadcastToRoom(executionId, map[string]interface{}{
		"name": "status",
		"data": map[string]interface{}{
			"status":       "waiting",
			"participants": nbPeoples, // L'hôte n'est pas compté
		},
	})
}

// handleNextQuestionEvent passe à la question suivante du quiz
// @Summary Passer à la question suivante
// @Description Envoie la prochaine question et ses réponses aux participants
// @Tags WebSocket
// @Accept json
// @Produce json
// @Param event body object true "Événement WebSocket 'nextQuestion'"
// @Success 200 {object} map[string]interface{} "Question et réponses envoyées aux participants"
// @Router /quiz/ws [post]
// @Security BearerAuth

func (sc *SocketController) handleNextQuestionEvent(data map[string]any) {
	executionId := data["executionId"].(string)
	quiz, err := sc.Service.QuizFromCode(executionId)
	if err != nil {
		return
	}

	if len(quiz.Questions) == 0 {
		return
	}

	nbPeoples, _ := sc.Service.GetRoomPeople(executionId)

	sc.broadcastToRoom(executionId, map[string]interface{}{
		"name": "status",
		"data": map[string]interface{}{
			"status":       "started",
			"participants": nbPeoples, // Toujours sans l'hôte
		},
	})

	sc.questionIdxMu.Lock()
	index := sc.questionIdx[executionId]
	sc.questionIdx[executionId]++
	sc.questionIdxMu.Unlock()

	if index >= len(quiz.Questions) {
		return // Toutes les questions ont été posées
	}

	question := quiz.Questions[index]
	var answers []string
	for _, answer := range question.Answers {
		answers = append(answers, answer.Title)
	}

	sc.broadcastToRoom(executionId, map[string]interface{}{
		"name": "newQuestion",
		"data": map[string]interface{}{
			"question": question.Title,
			"answers":  answers,
		},
	})
}

func (sc *SocketController) broadcastToRoom(executionId string, message map[string]interface{}) {
	sc.roomsMu.Lock()
	defer sc.roomsMu.Unlock()

	res, _ := json.Marshal(message)

	for _, conn := range sc.rooms[executionId] {
		_ = conn.WriteMessage(websocket.TextMessage, res)
	}

	if host, exists := sc.hosts[executionId]; exists {
		_ = host.WriteMessage(websocket.TextMessage, res)
	}
}
