package users

import (
	"cloud.google.com/go/firestore"
	"context"
)

type fireStoreAdapter struct {
	client *firestore.Client
}

func ConfigureStore(client *firestore.Client) Store {
	return &fireStoreAdapter{client}
}

func (fs *fireStoreAdapter) Upsert(user Document) error {
	_, err := fs.client.
		Collection("users").
		Doc(user.Uid).
		Set(context.Background(), user)
	return err
}

func (fs *fireStoreAdapter) GetUnique(uid string) (Document, error) {
	doc, err := fs.client.
		Collection("users").
		Doc(uid).
		Get(context.Background())

	if err != nil {
		return Document{}, err
	}

	if !doc.Exists() {
		return Document{}, ErrNotFound
	}

	var user Document
	if err2 := doc.DataTo(&user); err2 != nil {
		return user, err2
	}

	return user, nil
}
