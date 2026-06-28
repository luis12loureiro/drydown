package user

import "time"

// Service is the business-logic port for users.
type Service interface {
	List() ([]User, error)
	Get(uuid string) (User, error)
	Create(u User) (User, error)
	Update(uuid string, u User) (User, error)
	Delete(uuid string) error
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

func (s *service) Get(uuid string) (User, error) {
	return s.repo.FindByUUID(uuid)
}

func (s *service) Create(u User) (User, error) {
	if !u.Valid() {
		return User{}, ErrInvalid
	}
	u.CreatedAt = time.Now()
	return s.repo.Create(u)
}

func (s *service) Update(uuid string, u User) (User, error) {
	if !u.Valid() {
		return User{}, ErrInvalid
	}
	existing, err := s.repo.FindByUUID(uuid)
	if err != nil {
		return User{}, err
	}
	u.UUID = uuid
	u.CreatedAt = existing.CreatedAt
	return s.repo.Update(u)
}

func (s *service) Delete(uuid string) error {
	return s.repo.Delete(uuid)
}
