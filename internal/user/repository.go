package user

// Repository is the persistence port for users.
type Repository interface {
	FindAll() ([]User, error)
	FindByID(id int) (User, error)
	Create(u User) (User, error)
	Update(u User) (User, error)
	Delete(id int) error
}
