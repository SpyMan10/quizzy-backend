package quiz

import (
	"errors"
)

var (
	ErrNotFound             = errors.New("user not found")
	ErrInvalidPatchOperator = errors.New("invalid patch operator")
	ErrInvalidPatchField    = errors.New("invalid patch field")
)

type FieldPatchOp struct {
	Op    string `json:"op"`
	Path  string `json:"path"`
	Value any    `json:"value"`
}
type Links struct {
	Create string `json:"create,omitempty"`
	Start  string `json:"start,omitempty"`
}

// Quiz describe available data for a quiz.
type Quiz struct {
	Id          string     `firestore:"-" json:"id"`
	Title       string     `firestore:"title" json:"title"`
	Description string     `firestore:"description" json:"description"`
	Questions   []Question `firestore:"-" json:"questions"`
	Links       Links      `firestore:"-" json:"_links,omitempty"`
	Code        string     `firestore:"code" json:"code,omitempty"`
}

type Question struct {
	Id      string   `firestore:"-" json:"id"`
	Title   string   `firestore:"title" json:"title"`
	Answers []Answer `firestore:"-" json:"answers"`
}

type Answer struct {
	Id        string `firestore:"-" json:"id"`
	Title     string `firestore:"title" json:"title"`
	IsCorrect bool   `firestore:"isCorrect" json:"isCorrect"`
}

type Store interface {
	// Upsert Store or update the given user, if no user with the given id exists,
	// it will be created, otherwise it will be updated.
	Upsert(ownerId string, quiz Quiz) error

	// GetUnique returns the user matching to the given uid,
	// otherwise ErrNotFound is returned.
	GetUnique(ownerId, uid string) (Quiz, error)

	// GetQuizzes returns all quiz owned by the given user.
	GetQuizzes(ownerId string) ([]Quiz, error)

	// Patch update the given quiz.
	Patch(ownerId, uid string, fields []FieldPatchOp) error

	GetUniqueQuestion(ownerId, quizId, questionId string) (Question, error)

	UpsertQuestion(ownerId, quizId string, question Question) error

	// UpdateQuestion patch the given
	UpdateQuestion(ownerId, quizId string, question Question) error
}

type QuizCodeResolver interface {
	BindCode(quiz Quiz) error
	UnbindCode(code string) error
	GetQuiz(code string) (string, error)
}
