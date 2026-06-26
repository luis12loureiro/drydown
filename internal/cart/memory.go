package cart

import "sync"

// MemoryRepository is an in-memory adapter implementing Repository. Carts live
// only for the lifetime of the process; swap for sqlite.go later.
type MemoryRepository struct {
	mu    sync.RWMutex
	carts map[string]Cart
}

// NewMemoryRepository returns an empty in-memory cart store.
func NewMemoryRepository() *MemoryRepository {
	return &MemoryRepository{carts: make(map[string]Cart)}
}

func (r *MemoryRepository) Find(id string) (Cart, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	c, ok := r.carts[id]
	if !ok {
		return Cart{}, ErrNotFound
	}
	return c, nil
}

func (r *MemoryRepository) Save(c Cart) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.carts[c.ID] = c
	return nil
}

func (r *MemoryRepository) Delete(id string) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	delete(r.carts, id)
	return nil
}
