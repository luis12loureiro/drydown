package product

import "sync"

// MemoryRepository is an in-memory adapter implementing Repository. It is a
// stand-in for the eventual SQLite adapter (sqlite.go) so the app runs without
// a database during early development.
type MemoryRepository struct {
	mu       sync.RWMutex
	products map[int]Product
	nextID   int
}

// NewMemoryRepository returns a repository pre-loaded with the seed catalogue.
func NewMemoryRepository() Repository {
	r := &MemoryRepository{
		products: make(map[int]Product),
		nextID:   1,
	}
	for _, p := range seed() {
		_, _ = r.Create(p)
	}
	return r
}

func (r *MemoryRepository) FindAll() ([]Product, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	return r.sorted(func(Product) bool { return true }), nil
}

func (r *MemoryRepository) FindByFamily(fam string) ([]Product, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	return r.sorted(func(p Product) bool { return p.Fam == fam }), nil
}

func (r *MemoryRepository) FindByID(id int) (Product, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	p, ok := r.products[id]
	if !ok {
		return Product{}, ErrNotFound
	}
	return p, nil
}

func (r *MemoryRepository) Create(p Product) (Product, error) {
	r.mu.Lock()
	defer r.mu.Unlock()
	p.ID = r.nextID
	r.nextID++
	r.products[p.ID] = p
	return p, nil
}

func (r *MemoryRepository) Update(p Product) (Product, error) {
	r.mu.Lock()
	defer r.mu.Unlock()
	if _, ok := r.products[p.ID]; !ok {
		return Product{}, ErrNotFound
	}
	r.products[p.ID] = p
	return p, nil
}

func (r *MemoryRepository) Delete(id int) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	if _, ok := r.products[id]; !ok {
		return ErrNotFound
	}
	delete(r.products, id)
	return nil
}

// sorted returns the products matching keep, ordered by id for stable output.
// The caller must hold at least a read lock.
func (r *MemoryRepository) sorted(keep func(Product) bool) []Product {
	out := make([]Product, 0, len(r.products))
	for id := 1; id < r.nextID; id++ {
		p, ok := r.products[id]
		if ok && keep(p) {
			out = append(out, p)
		}
	}
	return out
}

// seed returns the starter catalogue, mirroring the original landing-page demo.
func seed() []Product {
	return []Product{
		{Brand: "Maison Noir", Name: "Oud Royale", Family: "Woody", Fam: "woody", Rating: 4.8, Reviews: 214, Badge: "Bestseller", Notes: []string{"Oud", "Saffron", "Rose"}, Sizes: []Size{{2, 16}, {5, 34}, {10, 58}}},
		{Brand: "Atelier Lumière", Name: "Fleur de Minuit", Family: "Floral", Fam: "floral", Rating: 4.6, Reviews: 132, Badge: "New", Notes: []string{"Tuberose", "Jasmine", "Musk"}, Sizes: []Size{{2, 13}, {5, 27}, {10, 46}}},
		{Brand: "Botanica", Name: "Vert Sauvage", Family: "Fresh", Fam: "fresh", Rating: 4.5, Reviews: 98, Notes: []string{"Bergamot", "Vetiver", "Mint"}, Sizes: []Size{{2, 11}, {5, 22}, {10, 38}}},
		{Brand: "Maison Noir", Name: "Ambre Solaire", Family: "Amber", Fam: "amber", Rating: 4.9, Reviews: 301, Badge: "Bestseller", Notes: []string{"Amber", "Vanilla", "Tonka"}, Sizes: []Size{{2, 18}, {5, 38}, {10, 64}}},
		{Brand: "Côte Marine", Name: "Sel & Soleil", Family: "Fresh", Fam: "fresh", Rating: 4.4, Reviews: 76, Notes: []string{"Sea Salt", "Bergamot", "Driftwood"}, Sizes: []Size{{2, 12}, {5, 24}, {10, 42}}},
		{Brand: "Atelier Lumière", Name: "Rose Noire", Family: "Floral", Fam: "floral", Rating: 4.7, Reviews: 187, Badge: "Limited", Notes: []string{"Rose", "Patchouli", "Oud"}, Sizes: []Size{{2, 17}, {5, 36}, {10, 60}}},
		{Brand: "Sucre Brûlé", Name: "Caramel Tabac", Family: "Gourmand", Fam: "gourmand", Rating: 4.6, Reviews: 142, Notes: []string{"Tonka", "Tobacco", "Caramel"}, Sizes: []Size{{2, 15}, {5, 31}, {10, 53}}},
		{Brand: "Maison Noir", Name: "Cuir Velvet", Family: "Leather", Fam: "leather", Rating: 4.8, Reviews: 166, Badge: "New", Notes: []string{"Leather", "Iris", "Cedar"}, Sizes: []Size{{2, 19}, {5, 39}, {10, 66}}},
	}
}
