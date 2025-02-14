package quiz

import "errors"

var (
	ErrNotFound             = errors.New("user not found")
	ErrInvalidPatchOperator = errors.New("invalid patch operator")
	ErrInvalidPatchField    = errors.New("invalid patch field")
)

type FieldPatchOp struct {
	Op    string      `json:"op"`
	Path  string      `json:"path"`
	Value interface{} `json:"value"`
}

// Document describe available data for a quiz.
type Document struct {
	Uid         string     `firestore:"uid" json:"id"`
	Title       string     `firestore:"title" json:"title"`
	Description string     `firestore:"description" json:"description"`
	Questions   []Question `firestore:"questions" json:"questions,omitempty"`
}

type Question struct {
	Uid     string   `firestore:"uid" json:"uid"`
	Title   string   `firestore:"title" json:"title"`
	Answers []Answer `firestore:"answers" json:"answers,omitempty"`
}

type Answer struct {
	Title     string `firestore:"title" json:"title"`
	IsCorrect bool   `firestore:"isCorrect" json:"isCorrect"`
}

type Store interface {
	// Upsert Store or update the given user, if no user with the given id exists,
	// it will be created, otherwise it will be updated.
	Upsert(ownerId string, quiz Document) error

	// GetUnique returns the user matching to the given uid,
	// otherwise ErrNotFound is returned.
	GetUnique(ownerId, uid string) (Document, error)

	// GetQuizzes returns all quiz owned by the given user.
	GetQuizzes(ownerId string) ([]Document, error)

	// Patch update the given quiz.
	Patch(ownerId, uid string, fields []FieldPatchOp) error
}
