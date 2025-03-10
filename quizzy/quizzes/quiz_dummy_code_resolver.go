package quizzes

import "fmt"

type dummyCodeResolver struct {
	entries map[string]string
	rooms   map[string]int
}

func (d *dummyCodeResolver) IncrRoomPeople(roomId string) error {
	v := d.rooms[roomId]
	d.rooms[roomId] = v + 1
	return nil
}

func (d *dummyCodeResolver) GetRoomPeople(roomId string) (int, error) {
	return d.rooms[roomId], nil
}

func (d *dummyCodeResolver) ResetRoomPeople(roomId string) error {
	d.rooms[roomId] = 0
	return nil
}

func (d *dummyCodeResolver) BindCode(ownerId string, quiz Quiz) error {
	d.entries[quiz.Code] = fmt.Sprintf("%s@%s", ownerId, quiz.Id)
	return nil
}

func (d *dummyCodeResolver) UnbindCode(code string) error {
	delete(d.entries, code)
	return nil
}

func (d *dummyCodeResolver) GetQuiz(code string) (string, error) {
	if quiz, ok := d.entries[code]; ok {
		return quiz, nil
	}

	return "", ErrNotFound
}
