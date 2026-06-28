package product

import "errors"

// ErrInvalid is returned when a product fails validation.
var ErrInvalid = errors.New("product: invalid")

// Service is the business-logic port for the catalogue.
type Service interface {
	List(family string) ([]Product, error)
	Get(uuid string) (Product, error)
	Create(p Product) (Product, error)
	Update(uuid string, p Product) (Product, error)
	Delete(uuid string) error
}

type service struct {
	repo Repository
}

// NewService wires a Service to its persistence adapter.
func NewService(repo Repository) Service {
	return &service{repo: repo}
}

func (s *service) List(family string) ([]Product, error) {
	if family == "" || family == "all" {
		return s.repo.FindAll()
	}
	return s.repo.FindByFamily(family)
}

func (s *service) Get(uuid string) (Product, error) {
	return s.repo.FindByUUID(uuid)
}

func (s *service) Create(p Product) (Product, error) {
	if !p.Valid() {
		return Product{}, ErrInvalid
	}
	return s.repo.Create(p)
}

func (s *service) Update(uuid string, p Product) (Product, error) {
	if !p.Valid() {
		return Product{}, ErrInvalid
	}
	p.UUID = uuid
	return s.repo.Update(p)
}

func (s *service) Delete(uuid string) error {
	return s.repo.Delete(uuid)
}
