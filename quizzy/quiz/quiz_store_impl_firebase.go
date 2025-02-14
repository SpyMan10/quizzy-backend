package quiz

import (
	"cloud.google.com/go/firestore"
	"context"
	"strings"
)

type fireStoreAdapter struct {
	client *firestore.Client
}

func ConfigureStore(client *firestore.Client) Store {
	return &fireStoreAdapter{client: client}
}

func (fs *fireStoreAdapter) Upsert(ownerId string, quiz Document) error {
	_, err := fs.client.
		Collection("users").
		Doc(ownerId).
		Collection("quizzes").
		Doc(quiz.Uid).
		Set(context.Background(), quiz)
	return err
}

func (fs *fireStoreAdapter) GetUnique(ownerId, uid string) (Document, error) {
	doc, err := fs.client.
		Collection("users").
		Doc(ownerId).
		Collection("quizzes").
		Doc(uid).
		Get(context.Background())

	if err != nil {
		return Document{}, err
	}

	if !doc.Exists() {
		return Document{}, ErrNotFound
	}

	var quiz Document
	if err2 := doc.DataTo(&quiz); err2 != nil {
		return quiz, err2
	}

	return quiz, nil
}

func (fs *fireStoreAdapter) GetQuizzes(ownerId string) ([]Document, error) {
	docsIter, err := fs.client.
		Collection("users").
		Doc(ownerId).Collection("quizzes").
		Documents(context.Background()).
		GetAll()

	if err != nil {
		return nil, err
	}

	// Must always be initialized to avoid nil pointer.
	arr := []Document{}
	for _, doc := range docsIter {
		quiz := Document{}
		if err2 := doc.DataTo(&quiz); err2 != nil {
			return nil, err2
		}

		arr = append(arr, quiz)
	}

	return arr, nil
}

func (fs *fireStoreAdapter) Patch(ownerId, uid string, fields []FieldPatchOp) error {
	var updates []firestore.Update
	for _, op := range fields {
		if op.Op != "replace" {
			return ErrInvalidPatchOperator
		}

		// Removing unwanted leading '/'.
		// Field path for firestore must not contain anyone of : ~/*[]
		path := strings.TrimLeft(op.Path, "/")

		if strings.ContainsAny(path, "*[]~") {
			return ErrInvalidPatchField
		}

		updates = append(updates, firestore.Update{
			FieldPath: strings.Split(path, "/"),
			Value:     op.Value,
		})
	}

	_, err := fs.client.
		Collection("users").
		Doc(ownerId).
		Collection("quizzes").
		Doc(uid).
		Update(context.Background(), updates)
	return err
}
