package user

import "time"

// Service is the business-logic port for users.
type Service interface {
	List() ([]User, error)
	Get(id int) (User, error)
	Create(u User) (User, error)
	Update(id int, u User) (User, error)
	Delete(id int) error
}

type service struct {
	repo Repository
}

// NewService wires a user Service to its persistence adapter.
func NewService(repo Repository) Service {
	return &service{repo: repo}
}

func (s *service) List() ([]User, error) {
	return s.repo.FindAll()
}

func (s *service) Get(id int) (User, error) {
	return s.repo.FindByID(id)
}

func (s *service) Create(u User) (User, error) {
	if !u.Valid() {
		return User{}, ErrInvalid
	}
	u.CreatedAt = time.Now()
	return s.repo.Create(u)
}

func (s *service) Update(id int, u User) (User, error) {
	if !u.Valid() {
		return User{}, ErrInvalid
	}
	existing, err := s.repo.FindByID(id)
	if err != nil {
		return User{}, err
	}
	u.ID = id
	u.CreatedAt = existing.CreatedAt // preserve immutable fields
	return s.repo.Update(u)
}

func (s *service) Delete(id int) error {
	return s.repo.Delete(id)
}
