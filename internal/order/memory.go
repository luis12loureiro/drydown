package order

import "sync"

// MemoryRepository is an in-memory adapter implementing Repository. Swap for
// sqlite.go later.
type MemoryRepository struct {
	mu     sync.RWMutex
	orders map[int]Order
	nextID int
}

// NewMemoryRepository returns an empty in-memory order store.
func NewMemoryRepository() *MemoryRepository {
	return &MemoryRepository{
		orders: make(map[int]Order),
		nextID: 1,
	}
}

func (r *MemoryRepository) FindAll() ([]Order, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	out := make([]Order, 0, len(r.orders))
	for id := 1; id < r.nextID; id++ {
		if o, ok := r.orders[id]; ok {
			out = append(out, o)
		}
	}
	return out, nil
}

func (r *MemoryRepository) FindByID(id int) (Order, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	o, ok := r.orders[id]
	if !ok {
		return Order{}, ErrNotFound
	}
	return o, nil
}

func (r *MemoryRepository) Create(o Order) (Order, error) {
	r.mu.Lock()
	defer r.mu.Unlock()
	o.ID = r.nextID
	r.nextID++
	r.orders[o.ID] = o
	return o, nil
}

func (r *MemoryRepository) Update(o Order) (Order, error) {
	r.mu.Lock()
	defer r.mu.Unlock()
	if _, ok := r.orders[o.ID]; !ok {
		return Order{}, ErrNotFound
	}
	r.orders[o.ID] = o
	return o, nil
}

func (r *MemoryRepository) Delete(id int) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	if _, ok := r.orders[id]; !ok {
		return ErrNotFound
	}
	delete(r.orders, id)
	return nil
}
