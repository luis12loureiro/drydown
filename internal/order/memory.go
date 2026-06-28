package order

import (
	"sync"

	"github.com/google/uuid"
)

// MemoryRepository is an in-memory adapter implementing Repository. Swap for
// sqlite.go later.
type memoryRepository struct {
	mu     sync.RWMutex
	orders map[string]Order
}

// NewMemoryRepository returns an empty in-memory order store.
func NewMemoryRepository() Repository {
	return &memoryRepository{
		orders: make(map[string]Order),
	}
}

func (r *memoryRepository) FindAll() ([]Order, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	out := make([]Order, 0, len(r.orders))
	for _, o := range r.orders {
		out = append(out, o)
	}
	return out, nil
}

func (r *memoryRepository) FindByUUID(id string) (Order, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	o, ok := r.orders[id]
	if !ok {
		return Order{}, ErrNotFound
	}
	return o, nil
}

func (r *memoryRepository) Create(o Order) (Order, error) {
	r.mu.Lock()
	defer r.mu.Unlock()
	o.UUID = uuid.New().String()
	r.orders[o.UUID] = o
	return o, nil
}

func (r *memoryRepository) Update(o Order) (Order, error) {
	r.mu.Lock()
	defer r.mu.Unlock()
	if _, ok := r.orders[o.UUID]; !ok {
		return Order{}, ErrNotFound
	}
	r.orders[o.UUID] = o
	return o, nil
}

func (r *memoryRepository) Delete(id string) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	if _, ok := r.orders[id]; !ok {
		return ErrNotFound
	}
	delete(r.orders, id)
	return nil
}
