package users

type DummyUserStoreImpl struct {
	Users []User
}

func NewDummyStore() Store {
	return &DummyUserStoreImpl{
		Users: make([]User, 0),
	}
}

func (st *DummyUserStoreImpl) Upsert(user User) error {
	for _, u := range st.Users {
		if u.Id == user.Id {
			u.Email = user.Email
			u.Username = user.Username
			return nil
		}
	}

	st.Users = append(st.Users, user)
	return nil
}

func (st *DummyUserStoreImpl) GetUnique(id string) (User, error) {
	for _, user := range st.Users {
		if user.Id == id {
			return user, nil
		}
	}

	return User{}, ErrNotFound
}
