package users

type dummyUserStoreImpl struct {
	Users []User
}

func _newDummyStore() Store {
	return &dummyUserStoreImpl{
		Users: make([]User, 0),
	}
}

func (st *dummyUserStoreImpl) Upsert(user User) error {
	for i, u := range st.Users {
		if u.Id == user.Id {
			st.Users[i] = user
		}
	}

	st.Users = append(st.Users, user)
	return nil
}

func (st *dummyUserStoreImpl) GetUnique(id string) (User, error) {
	for _, user := range st.Users {
		if user.Id == id {
			return user, nil
		}
	}

	return User{}, ErrNotFound
}
