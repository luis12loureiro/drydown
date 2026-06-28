package product

import (
	"sync"

	"github.com/google/uuid"
)

// MemoryRepository is an in-memory adapter implementing Repository. It is a
// stand-in for the eventual SQLite adapter (sqlite.go) so the app runs without
// a database during early development.
type memoryRepository struct {
	mu       sync.RWMutex
	products map[string]Product
}

// NewMemoryRepository returns a repository pre-loaded with the seed catalogue.
func NewMemoryRepository() Repository {
	r := &memoryRepository{
		products: make(map[string]Product),
	}
	for _, p := range seed() {
		_, _ = r.Create(p)
	}
	return r
}

func (r *memoryRepository) FindAll() ([]Product, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	out := make([]Product, 0, len(r.products))
	for _, p := range r.products {
		out = append(out, p)
	}
	return out, nil
}

func (r *memoryRepository) FindByFamily(fam string) ([]Product, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	var out []Product
	for _, p := range r.products {
		if p.Fam == fam {
			out = append(out, p)
		}
	}
	return out, nil
}

func (r *memoryRepository) FindByUUID(id string) (Product, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	p, ok := r.products[id]
	if !ok {
		return Product{}, ErrNotFound
	}
	return p, nil
}

func (r *memoryRepository) Create(p Product) (Product, error) {
	r.mu.Lock()
	defer r.mu.Unlock()
	p.UUID = uuid.New().String()
	r.products[p.UUID] = p
	return p, nil
}

func (r *memoryRepository) Update(p Product) (Product, error) {
	r.mu.Lock()
	defer r.mu.Unlock()
	if _, ok := r.products[p.UUID]; !ok {
		return Product{}, ErrNotFound
	}
	r.products[p.UUID] = p
	return p, nil
}

func (r *memoryRepository) Delete(id string) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	if _, ok := r.products[id]; !ok {
		return ErrNotFound
	}
	delete(r.products, id)
	return nil
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
