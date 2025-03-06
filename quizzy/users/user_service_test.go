package users

import (
	"fmt"
	"github.com/google/uuid"
	"math/rand"
	"testing"
)

const userCount = 3

func _mockUserService(svc UserService) []string {
	ids := make([]string, userCount)

	for i := range ids {
		id := uuid.New().String()
		ids[i] = id
		_ = svc.Create(User{
			Id:       id,
			Username: fmt.Sprintf("test-user%d", i+1),
			Email:    fmt.Sprintf("test.user%d@mail.com", i+1),
		})
	}

	return ids
}

func _createDummyService() UserService {
	return &UserServiceImpl{Store: _newDummyStore()}
}

func TestCreateAndGetUser(t *testing.T) {
	s := _createDummyService()
	id := uuid.New().String()

	if e := s.Create(User{
		Id:       id,
		Username: "test-user",
		Email:    "test.user@mail.com",
	}); e != nil {
		t.Fatalf("create user: %s", e)
	}

	if u, e := s.Get(id); e != nil {
		t.Fatalf("get user: %s", e)
	} else if u.Id != id {
		t.Fatalf("get user: id %s != %s", u.Id, id)
	}
}

func TestUpdateUsername(t *testing.T) {
	s := _createDummyService()
	ids := _mockUserService(s)

	// Randomize user selection for update.
	if u, e := s.Get(ids[rand.Intn(userCount)]); e != nil {
		t.Fatalf("get user by id after init: %s", e)
	} else {
		t.Logf("got user with id %s (username: %s)", u.Id, u.Username)
		u.Username = "updated-username"
		if e2 := s.Update(u); e2 != nil {
			t.Fatalf("update username")
		}
		t.Logf("update done for user with id %s (username: %s", u.Id, u.Username)

		if u2, e2 := s.Get(u.Id); e2 != nil {
			t.Fatalf("get user by id after update: %s", e2)
		} else {
			t.Logf("got user with id %s (username: %s)", u2.Id, u2.Username)
			if u2.Username != "updated-username" {
				t.Fatalf("check if username was updated: (stored) %s != %s", u2.Username, u.Username)
			}
		}
	}
}
