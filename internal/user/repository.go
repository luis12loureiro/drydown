package user

// Repository is the persistence port for users.
type Repository interface {
	FindAll() ([]User, error)
	FindByUUID(uuid string) (User, error)
	Create(u User) (User, error)
	Update(u User) (User, error)
	Delete(uuid string) error
}
