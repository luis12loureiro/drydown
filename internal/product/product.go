// Package product holds the catalogue domain: the Product entity and its rules.
package product

import "errors"

// ErrNotFound is returned when a product cannot be located by its uuid.
var ErrNotFound = errors.New("product: not found")

// Size is a single decant volume offered for a product and its price (in euros).
type Size struct {
	ML    int `json:"ml"`
	Price int `json:"price"`
}

// Product is a fragrance available in the catalogue. Attributes are kept
// deliberately simple for the first deploy.
type Product struct {
	UUID    string
	Brand   string
	Name    string
	Family  string // display label, e.g. "Woody"
	Fam     string // gradient/filter key: woody|floral|fresh|amber|gourmand|leather
	Rating  float64
	Reviews int
	Badge   string // "" when none
	Notes   []string
	Sizes   []Size
}

// PriceFor returns the price of the given decant size and whether it exists.
func (p Product) PriceFor(ml int) (int, bool) {
	for _, s := range p.Sizes {
		if s.ML == ml {
			return s.Price, true
		}
	}
	return 0, false
}

// Valid reports whether the product carries the minimum data needed to be sold.
func (p Product) Valid() bool {
	return p.Name != "" && p.Brand != "" && len(p.Sizes) > 0
}
