package product

import "errors"

// ErrInvalid is returned when a product fails validation.
var ErrInvalid = errors.New("product: invalid")

// Service is the business-logic port for the catalogue.
type Service interface {
	List(family string) ([]Product, error)
	Get(id int) (Product, error)
	Create(p Product) (Product, error)
	Update(id int, p Product) (Product, error)
	Delete(id int) error
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

func (s *service) Get(id int) (Product, error) {
	return s.repo.FindByID(id)
}

func (s *service) Create(p Product) (Product, error) {
	if !p.Valid() {
		return Product{}, ErrInvalid
	}
	return s.repo.Create(p)
}

func (s *service) Update(id int, p Product) (Product, error) {
	if !p.Valid() {
		return Product{}, ErrInvalid
	}
	p.ID = id
	return s.repo.Update(p)
}

func (s *service) Delete(id int) error {
	return s.repo.Delete(id)
}
