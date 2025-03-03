package users

import (
	"cloud.google.com/go/firestore"
	"context"
	"strings"
)

type fireStoreAdapter struct {
	client *firestore.Client
}

func ConfigureStore(client *firestore.Client) Store {
	return &fireStoreAdapter{client}
}

func (fs *fireStoreAdapter) Upsert(user User) error {
	_, err := fs.client.
		Doc(strings.Join([]string{"users", user.Id}, "/")).
		Set(context.Background(), user)
	return err
}

func (fs *fireStoreAdapter) GetUnique(id string) (User, error) {
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
