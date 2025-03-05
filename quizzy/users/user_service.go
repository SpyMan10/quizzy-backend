package users

type UserService interface {
	// Get returns the matching User with the given unique id.
	Get(id string) (User, error)

	// Update will update the given User.
	Update(user User) error

	// Create the given User in our application.
	Create(user User) error
}
