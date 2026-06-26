package order

import "time"

// Service is the business-logic port for orders.
type Service interface {
	List() ([]Order, error)
	Get(id int) (Order, error)
	Place(items []Item) (Order, error)
	SetStatus(id int, status string) (Order, error)
	Cancel(id int) error
}

type service struct {
	repo Repository
}

// NewService wires an order Service to its persistence adapter.
func NewService(repo Repository) Service {
	return &service{repo: repo}
}

func (s *service) List() ([]Order, error) {
	return s.repo.FindAll()
}

func (s *service) Get(id int) (Order, error) {
	return s.repo.FindByID(id)
}

// Place creates a pending order from the supplied lines, computing the total.
func (s *service) Place(items []Item) (Order, error) {
	if len(items) == 0 {
		return Order{}, ErrEmpty
	}
	total := 0
	for _, it := range items {
		total += it.Price * it.Qty
	}
	return s.repo.Create(Order{
		Items:     items,
		Total:     total,
		Status:    StatusPending,
		CreatedAt: time.Now(),
	})
}

func (s *service) SetStatus(id int, status string) (Order, error) {
	o, err := s.repo.FindByID(id)
	if err != nil {
		return Order{}, err
	}
	o.Status = status
	return s.repo.Update(o)
}

func (s *service) Cancel(id int) error {
	return s.repo.Delete(id)
}
