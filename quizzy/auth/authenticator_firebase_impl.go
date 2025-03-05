package auth

import (
	"context"
	"quizzy.app/backend/quizzy/services"
)

type FirebaseAuthenticator struct {
	Fbs *services.FirebaseServices
}

func (auth *FirebaseAuthenticator) Authorize(token string) (Identity, error) {
	if tk, err := auth.Fbs.Auth.VerifyIDTokenAndCheckRevoked(context.Background(), token); err != nil {
		return Identity{}, err
	} else {
		return Identity{
			Token: token,
			Uid:   tk.UID,
			Email: tk.Claims["email"].(string),
		}, err
	}
}
