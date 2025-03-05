package users

type UserServiceImpl struct {
	Store Store
}

func (us *UserServiceImpl) Create(user User) error {
	return us.Store.Upsert(user)
}

func (us *UserServiceImpl) Update(user User) error {
	return us.Store.Upsert(user)
}

func (us *UserServiceImpl) Get(id string) (User, error) {
	return us.Store.GetUnique(id)
}
