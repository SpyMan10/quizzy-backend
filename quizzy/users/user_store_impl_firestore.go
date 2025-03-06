package users

import (
	"cloud.google.com/go/firestore"
	"context"
	"strings"
)

type userFirestore struct {
	client *firestore.Client
}

func NewFirestore(client *firestore.Client) Store {
	return &userFirestore{client}
}

func (fs *userFirestore) Upsert(user User) error {
	_, err := fs.client.
		Doc(strings.Join([]string{"users", user.Id}, "/")).
		Set(context.Background(), user)
	return err
}

func (fs *userFirestore) GetUnique(id string) (User, error) {
	doc, err := fs.client.
		Doc(strings.Join([]string{"users", id}, "/")).
		Get(context.Background())

	if err != nil {
		return User{}, err
	}

	if !doc.Exists() {
		return User{}, ErrNotFound
	}

	var user User
	if err2 := doc.DataTo(&user); err2 != nil {
		return user, err2
	}

	user.Id = doc.Ref.ID
	return user, nil
}
