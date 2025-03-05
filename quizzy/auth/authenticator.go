package auth

type Identity struct {
	Token string `json:"-"`
	Uid   string `json:"-"`
	Email string `json:"-"`
}

type Authenticator interface {
	Authorize(token string) (Identity, error)
}
