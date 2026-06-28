package user

import (
	"sync"

	"github.com/google/uuid"
)

// MemoryRepository is an in-memory adapter implementing Repository. Swap for
// sqlite.go later.
type memoryRepository struct {
	mu    sync.RWMutex
	users map[string]User
}

// NewMemoryRepository returns an empty in-memory user store.
func NewMemoryRepository() Repository {
	return &memoryRepository{
		users: make(map[string]User),
	}
}

func (r *memoryRepository) FindAll() ([]User, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	out := make([]User, 0, len(r.users))
	for _, u := range r.users {
		out = append(out, u)
	}
	return out, nil
}

func (r *memoryRepository) FindByUUID(id string) (User, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	u, ok := r.users[id]
	if !ok {
		return User{}, ErrNotFound
	}
	return u, nil
}

func (r *memoryRepository) Create(u User) (User, error) {
	r.mu.Lock()
	defer r.mu.Unlock()
	u.UUID = uuid.New().String()
	r.users[u.UUID] = u
	return u, nil
}

func (r *memoryRepository) Update(u User) (User, error) {
	r.mu.Lock()
	defer r.mu.Unlock()
	if _, ok := r.users[u.UUID]; !ok {
		return User{}, ErrNotFound
	}
	r.users[u.UUID] = u
	return u, nil
}

func (r *memoryRepository) Delete(id string) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	if _, ok := r.users[id]; !ok {
		return ErrNotFound
	}
	delete(r.users, id)
	return nil
}
