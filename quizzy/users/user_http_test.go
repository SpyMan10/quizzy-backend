package users

import (
	"github.com/gavv/httpexpect/v2"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"net/http"
	"quizzy.app/backend/quizzy/auth"
	"testing"
)

func _configureTestHandler(id auth.Identity, users []User) http.HandlerFunc {
	eng := gin.Default()
	rt := eng.Group("", auth.ProvideAuthenticator(&auth.DummyAuthenticator{PlaceHolder: id}))
	con := Controller{
		Service: &UserServiceImpl{Store: _newDummyStore(users)},
	}
	con.ConfigureRouting(rt)

	return eng.ServeHTTP
}

func _fakeId() auth.Identity {
	return auth.Identity{
		Token: "x",
		Uid:   uuid.New().String(),
		Email: "test@mail.net",
	}
}

func TestPostUser(t *testing.T) {
	id := _fakeId()
	ex := httpexpect.Default(t, "/")
	handler := _configureTestHandler(id, nil)

	ex.POST("/users").
		WithHeader("Authorization", "Bearer x").
		WithJSON(createUserRequest{Username: "dummy-user"}).
		WithHandler(handler).
		Expect().
		Status(http.StatusCreated).
		NoContent()

	var user User
	ex.GET("/users/me").
		WithHeader("Authorization", "Bearer x").
		WithHandler(handler).
		Expect().
		Status(http.StatusOK).
		JSON().Object().Decode(&user)

	assert.Equal(t, user.Id, id.Uid)
	assert.Equal(t, user.Email, id.Email)
	assert.Equal(t, user.Username, "dummy-user")
}

func TestPostUserWithoutAuthorization(t *testing.T) {
	id := _fakeId()
	ex := httpexpect.Default(t, "/")
	handler := _configureTestHandler(id, nil)

	ex.POST("/users").
		WithJSON(createUserRequest{Username: "dummy-user"}).
		WithHandler(handler).
		Expect().
		Status(http.StatusUnauthorized).
		NoContent()
}

func TestPostUserWithoutPayload(t *testing.T) {
	id := _fakeId()
	ex := httpexpect.Default(t, "/")
	handler := _configureTestHandler(id, nil)

	ex.POST("/users").
		WithHeader("Authorization", "Bearer x").
		WithHandler(handler).
		Expect().
		Status(http.StatusBadRequest).
		NoContent()
}

func TestGetUserSelf(t *testing.T) {
	id := _fakeId()
	ex := httpexpect.Default(t, "/")
	handler := _configureTestHandler(id, []User{
		{Username: "dummy-user", Email: id.Email, Id: id.Uid},
	})

	var user User
	ex.GET("/users/me").
		WithHeader("Authorization", "Bearer x").
		WithHandler(handler).
		Expect().
		Status(http.StatusOK).
		JSON().Object().Decode(&user)

	assert.Equal(t, user.Id, id.Uid)
	assert.Equal(t, user.Email, id.Email)
	assert.Equal(t, user.Username, "dummy-user")
}
