package users

import (
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"testing"
)

func _createDummyService(data []User) UserService {
	return &UserServiceImpl{Store: _newDummyStore(data)}
}

func TestUserServiceCreate(t *testing.T) {
	data := make([]User, 0)
	s := _createDummyService(data)

	expected := User{
		Id:       uuid.New().String(),
		Username: "test-user",
		Email:    "test.user@mail.com",
	}
	e := s.Create(expected)

	assert.Nil(t, e)
	assert.Equal(t, data[0].Id, expected.Id)
	assert.Equal(t, data[0].Username, expected.Username)
	assert.Equal(t, data[0].Email, expected.Email)
}

func TestUserServiceUpdateUsername(t *testing.T) {
	id := uuid.New().String()
	data := []User{
		{
			Id:       id,
			Username: "test-user",
			Email:    "test.user@mail.com",
		},
	}
	s := _createDummyService(data)

	usr := User{
		Id:       id,
		Username: "test-user-333",
		Email:    "test.user@mail.com",
	}

	e := s.Update(usr)
	assert.Nil(t, e)
	assert.Equal(t, usr.Id, data[0].Id)
	assert.Equal(t, usr.Username, data[0].Username)
	assert.Equal(t, usr.Email, data[0].Email)
}
