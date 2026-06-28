package cart

import "sync"

// MemoryRepository is an in-memory adapter implementing Repository. Carts live
// only for the lifetime of the process; swap for sqlite.go later.
type memoryRepository struct {
	mu    sync.RWMutex
	carts map[string]Cart
}

// NewMemoryRepository returns an empty in-memory cart store.
func NewMemoryRepository() Repository {
	return &memoryRepository{carts: make(map[string]Cart)}
}

func (r *memoryRepository) Find(id string) (Cart, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	c, ok := r.carts[id]
	if !ok {
		return Cart{}, ErrNotFound
	}
	return c, nil
}

func (r *memoryRepository) Save(c Cart) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.carts[c.UUID] = c
	return nil
}

func (r *memoryRepository) Delete(id string) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	delete(r.carts, id)
	return nil
}
