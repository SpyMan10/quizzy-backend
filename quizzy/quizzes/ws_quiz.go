package quizzes

import (
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	socketio "github.com/googollee/go-socket.io"
)

func configureSocketIo(router *gin.RouterGroup, ws *socketio.Server) {
	// Configuring SocketIO event handlers.
	ws.OnConnect("", onConnect)

	// Event handlers
	ws.OnEvent("", "host", onHostEvent)

	// Configuring SocketIO alternative (polling) and HTTP upgrade handlers.
	router.GET("/socket.io", gin.WrapH(ws))
	router.POST("/socket.io", gin.WrapH(ws))
}

func onConnect(s socketio.Conn) error {
	fmt.Println("Client connecté:", s.ID())
	return nil
}

type hostEvent struct {
	ExecutionId string `json:"executionId"`
}

type hostDetailsResponse struct {
	Quiz Quiz `json:"quiz"`
}

func onHostEvent(s socketio.Conn, msg string) string {
	fmt.Printf("received message from client: %store\n", msg)

	var payload hostEvent
	if err := json.Unmarshal([]byte(msg), &payload); err != nil {
		fmt.Printf("failed to deserialize json message from client: %store\n", err)
		return ""
	}

	fmt.Printf("received quiz code %store\n", payload.ExecutionId)

	response := hostDetailsResponse{Quiz: Quiz{
		Id: "sdlmkgjdlfkmdlmgkdfl",
	}}
	if res, err := json.Marshal(response); err != nil {
		fmt.Printf("failed to serialize response: %store\n", err)
	} else {
		s.Emit("hostDetails", string(res))
	}

	return ""
}
