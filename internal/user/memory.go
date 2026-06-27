package user

import "sync"

// MemoryRepository is an in-memory adapter implementing Repository. Swap for
// sqlite.go later.
type MemoryRepository struct {
	mu     sync.RWMutex
	users  map[int]User
	nextID int
}

// NewMemoryRepository returns an empty in-memory user store.
func NewMemoryRepository() Repository {
	return &MemoryRepository{
		users:  make(map[int]User),
		nextID: 1,
	}
}

func (r *MemoryRepository) FindAll() ([]User, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	out := make([]User, 0, len(r.users))
	for id := 1; id < r.nextID; id++ {
		if u, ok := r.users[id]; ok {
			out = append(out, u)
		}
	}
	return out, nil
}

func (r *MemoryRepository) FindByID(id int) (User, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	u, ok := r.users[id]
	if !ok {
		return User{}, ErrNotFound
	}
	return u, nil
}

func (r *MemoryRepository) Create(u User) (User, error) {
	r.mu.Lock()
	defer r.mu.Unlock()
	u.ID = r.nextID
	r.nextID++
	r.users[u.ID] = u
	return u, nil
}

func (r *MemoryRepository) Update(u User) (User, error) {
	r.mu.Lock()
	defer r.mu.Unlock()
	if _, ok := r.users[u.ID]; !ok {
		return User{}, ErrNotFound
	}
	r.users[u.ID] = u
	return u, nil
}

func (r *MemoryRepository) Delete(id int) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	if _, ok := r.users[id]; !ok {
		return ErrNotFound
	}
	delete(r.users, id)
	return nil
}
