package users

import "errors"

var (
	ErrNotFound = errors.New("user not found")
)

type Document struct {
	Uid      string `firestore:"uid" json:"uid"`
	Username string `firestore:"username" json:"username"`
	Email    string `firestore:"email" json:"email"`
}

type Store interface {
	// Upsert Store or update the given user, if no user with the given id exists,
	// it will be created, otherwise it will be updated.
	Upsert(user Document) error

	// GetUnique returns the user matching to the given uid,
	// otherwise ErrNotFound is returned.
	GetUnique(uid string) (Document, error)
}
