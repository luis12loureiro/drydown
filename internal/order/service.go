package order

import "time"

// Service is the business-logic port for orders.
type Service interface {
	List() ([]Order, error)
	Get(uuid string) (Order, error)
	Place(items []Item) (Order, error)
	SetStatus(uuid string, status string) (Order, error)
	Cancel(uuid string) error
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

func (s *service) Get(uuid string) (Order, error) {
	return s.repo.FindByUUID(uuid)
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

func (s *service) SetStatus(uuid string, status string) (Order, error) {
	o, err := s.repo.FindByUUID(uuid)
	if err != nil {
		return Order{}, err
	}
	o.Status = status
	return s.repo.Update(o)
}

func (s *service) Cancel(uuid string) error {
	return s.repo.Delete(uuid)
}
