package users

import (
	"cloud.google.com/go/firestore"
	"context"
)

type fireStoreAdapter struct {
	client *firestore.Client
}

func ConfigureUserStore(client *firestore.Client) UserStore {
	return &fireStoreAdapter{client}
}

func (fs *fireStoreAdapter) Upsert(user User) error {
	_, err := fs.client.Collection("users").Doc(user.Uid).Set(context.Background(), user)
	return err
}

func (fs *fireStoreAdapter) GetUnique(uid string) (User, error) {
	doc, err := fs.client.Collection("users").Doc(uid).Get(context.Background())
	if err != nil {
		return User{}, err
	}

	if !doc.Exists() {
		return User{}, ErrNotFound
	}

	return User{
		Uid:      uid,
		Username: doc.Data()["Username"].(string),
	}, nil
}
