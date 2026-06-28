package cart

import (
	"errors"

	"github.com/luis12loureiro/drydown/internal/product"
)

// ErrInvalidSize is returned when the requested decant size is not sold.
var ErrInvalidSize = errors.New("cart: invalid size")

// Catalog is the port the cart needs from the product domain to price lines.
// product.Service satisfies it.
type Catalog interface {
	Get(uuid string) (product.Product, error)
}

// Service is the business-logic port for the cart.
type Service interface {
	Get(id string) (Cart, error)
	AddItem(cartID string, productUUID string, ml, qty int) (Cart, error)
	RemoveItem(cartID string, productUUID string, ml int) (Cart, error)
	Clear(cartID string) error
}

type service struct {
	repo    Repository
	catalog Catalog
}

// NewService wires the cart to its store and the product catalogue.
func NewService(repo Repository, catalog Catalog) Service {
	return &service{repo: repo, catalog: catalog}
}

// Get returns the visitor's cart, or a fresh empty one when none exists yet.
func (s *service) Get(id string) (Cart, error) {
	c, err := s.repo.Find(id)
	if errors.Is(err, ErrNotFound) {
		return Cart{UUID: id}, nil
	}
	return c, err
}

func (s *service) AddItem(cartID string, productUUID string, ml, qty int) (Cart, error) {
	if qty < 1 {
		qty = 1
	}
	c, err := s.Get(cartID)
	if err != nil {
		return Cart{}, err
	}
	p, err := s.catalog.Get(productUUID)
	if err != nil {
		return Cart{}, err
	}
	price, ok := p.PriceFor(ml)
	if !ok {
		return Cart{}, ErrInvalidSize
	}
	c.Add(Item{
		ProductUUID: p.UUID,
		Brand:       p.Brand,
		Name:        p.Name,
		Fam:         p.Fam,
		ML:          ml,
		Price:       price,
		Qty:         qty,
	})
	if err := s.repo.Save(c); err != nil {
		return Cart{}, err
	}
	return c, nil
}

func (s *service) RemoveItem(cartID string, productUUID string, ml int) (Cart, error) {
	c, err := s.Get(cartID)
	if err != nil {
		return Cart{}, err
	}
	c.Remove(productUUID, ml)
	if err := s.repo.Save(c); err != nil {
		return Cart{}, err
	}
	return c, nil
}

func (s *service) Clear(cartID string) error {
	return s.repo.Delete(cartID)
}
