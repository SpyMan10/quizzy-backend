package services

import (
	"cloud.google.com/go/firestore"
	"context"
	"errors"
	firebase "firebase.google.com/go"
	fireauth "firebase.google.com/go/auth"
	"google.golang.org/api/option"
	"quizzy.app/backend/quizzy/cfg"
)

var (
	ErrFirebaseConfNotFound = errors.New("no firebase config file found: please set APP_FIREBASE_CONF_FILE to valid firebase configuration file")
)

type FirebaseServices struct {
	Store *firestore.Client
	Auth  *fireauth.Client
}

func ConfigureFirebase(cfg cfg.AppConfig) (FirebaseServices, error) {
	if len(cfg.FirebaseConfFile) == 0 {
		return FirebaseServices{}, ErrFirebaseConfNotFound
	}

	opt := option.WithCredentialsFile(cfg.FirebaseConfFile)
	if app, err := firebase.NewApp(context.Background(), nil, opt); app != nil && err == nil {
		store, _ := app.Firestore(context.Background())
		auth, _ := app.Auth(context.Background())
		return FirebaseServices{
			Store: store,
			Auth:  auth,
		}, nil
	} else {
		return FirebaseServices{}, err
	}
}
