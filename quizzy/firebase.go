package quizzy

import (
	"cloud.google.com/go/firestore"
	"context"
	"errors"
	firebase "firebase.google.com/go"
	"google.golang.org/api/option"
)

var (
	ErrFirebaseConfNotFound = errors.New("no firebase config file found: please set APP_FIREBASE_CONF_FILE to valid firebase configuration file")
)

func ConfigureFirebase(cfg AppConfig) (*firestore.Client, error) {
	if len(cfg.firebaseConfFile) == 0 {
		return nil, ErrFirebaseConfNotFound
	}

	opt := option.WithCredentialsFile(cfg.firebaseConfFile)
	if app, err := firebase.NewApp(context.Background(), nil, opt); app != nil && err == nil {
		return app.Firestore(context.Background())
	} else {
		return nil, err
	}
}
