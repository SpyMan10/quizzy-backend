package quizzes

import "fmt"

type dummyCodeResolver struct {
	entries map[string]string
}

func (d dummyCodeResolver) BindCode(ownerId string, quiz Quiz) error {
	d.entries[quiz.Code] = fmt.Sprintf("%s@%s", ownerId, quiz.Id)
	return nil
}

func (d dummyCodeResolver) UnbindCode(code string) error {
	delete(d.entries, code)
	return nil
}

func (d dummyCodeResolver) GetQuiz(code string) (string, error) {
	if quiz, ok := d.entries[code]; ok {
		return quiz, nil
	}

	return "", ErrNotFound
}
