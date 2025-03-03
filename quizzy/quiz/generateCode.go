package quiz

import (
	"crypto/rand"
	"fmt"
)

const CodeLength = 3

func GenerateCode() (string, error) {
	b := make([]byte, CodeLength)
	if _, err := rand.Read(b); err != nil {
		return "", err
	}
	return fmt.Sprintf("%X", b), nil
}
