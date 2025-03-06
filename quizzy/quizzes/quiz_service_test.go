package quizzes

import (
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"testing"
)

func _createDummyQuizService() QuizService {
	return &QuizServiceImpl{
		store:    &dummyQuizStoreImpl{entries: make([]dummyEntry, 0)},
		resolver: &dummyCodeResolver{entries: make(map[string]string)},
	}
}

func TestCreateAndGetQuiz(t *testing.T) {
	svc := _createDummyQuizService()
	ownerId := uuid.New().String()

	code, err := GenerateCode()
	if err != nil {
		t.Fatalf("failed to generate unique code: %s", err)
	}

	expected := Quiz{
		Id:          uuid.New().String(),
		Title:       "quiz-title",
		Description: "test-description",
		Questions:   make([]Question, 0),
		Code:        code,
	}

	err2 := svc.Create(ownerId, expected)

	if err2 != nil {
		t.Fatalf("failed to create quiz: %s", err2)
	}

	// Try to get newly created quiz.

	if quiz, err3 := svc.Get(ownerId, expected.Id); err3 != nil {
		t.Fatalf("failed to get quiz: %s", err3)
	} else {
		assert.Equal(t, quiz.Title, expected.Title)
		assert.Equal(t, quiz.Description, expected.Description)
		assert.Equal(t, len(quiz.Questions), len(expected.Questions))
		assert.Equal(t, quiz.Code, expected.Code)
	}
}
