package auth

type DummyAuthenticator struct {
	PlaceHolder Identity
}

func (d *DummyAuthenticator) Authorize(token string) (Identity, error) {
	return d.PlaceHolder, nil
}
