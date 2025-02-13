package users

import "errors"

var (
	ErrNotFound = errors.New("user not found")
)

type User struct {
	Uid      string `json:"uid"`
	Username string `json:"username"`
}

type UserStore interface {
	// Upsert Store or update the given user, if no user with the given id exists,
	// it will be created, otherwise it will be updated.
	Upsert(user User) error

	// GetUnique returns the user matching to the given uid, otherwise ErrNotFound is returned.
	GetUnique(uid string) (User, error)
}
