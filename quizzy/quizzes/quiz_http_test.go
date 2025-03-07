package quizzes

import (
	"fmt"
	"github.com/gavv/httpexpect/v2"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"net/http"
	"quizzy.app/backend/quizzy/auth"
	"testing"
)

func _configureTestHandler(id auth.Identity, data []dummyEntry) http.HandlerFunc {
	eng := gin.Default()
	rt := eng.Group("", auth.ProvideAuthenticator(&auth.DummyAuthenticator{PlaceHolder: id}))
	con := Controller{
		Service: &QuizServiceImpl{store: _newDummyStore(data)},
	}
	con.ConfigureRouting(rt)

	return eng.ServeHTTP
}

func _fakeId() auth.Identity {
	return auth.Identity{
		Token: "x",
		Uid:   uuid.New().String(),
		Email: "test@mail.net",
	}
}

func TestPostQuiz(t *testing.T) {
	id := _fakeId()
	handler := _configureTestHandler(id, nil)
	ex := httpexpect.Default(t, "")

	var quiz Quiz
	payload := CreateQuizRequest{Title: "test-quiz", Description: "desc-test-quiz"}
	resp := ex.POST("/quiz").
		WithHandler(handler).
		WithHeader("Authorization", "Bearer x").
		WithJSON(&payload).
		Expect().
		Status(http.StatusCreated)

	resp.JSON().Object().Decode(&quiz)

	assert.Equal(t, quiz.Title, payload.Title)
	assert.Equal(t, quiz.Description, payload.Description)

	resp.Headers().
		Value("Location").
		Array().
		HasValue(0, fmt.Sprintf("http://localhost:8000/quiz/%s", quiz.Id))
}

func TestPostQuizWithoutAuthorization(t *testing.T) {
	handler := _configureTestHandler(_fakeId(), nil)
	ex := httpexpect.Default(t, "")
	
	ex.POST("/quiz").
		WithHandler(handler).
		Expect().
		NoContent().
		Status(http.StatusUnauthorized)
}
