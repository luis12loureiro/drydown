// Package cart holds the shopping-cart domain: a Cart and the lines within it.
package cart

import "errors"

// ErrNotFound is returned when no cart exists for a given id.
var ErrNotFound = errors.New("cart: not found")

// Item is a single line in the cart: a product at a chosen decant size.
type Item struct {
	ProductUUID string
	Brand       string
	Name        string
	Fam         string
	ML          int
	Price       int // unit price for the chosen size, in euros
	Qty         int
}

// Subtotal is the line total in euros.
func (i Item) Subtotal() int {
	return i.Price * i.Qty
}

// Cart is a visitor's bag, identified by an opaque session uuid.
type Cart struct {
	UUID  string
	Items []Item
}

// Count is the total number of decants across all lines.
func (c Cart) Count() int {
	n := 0
	for _, it := range c.Items {
		n += it.Qty
	}
	return n
}

// Total is the cart subtotal in euros.
func (c Cart) Total() int {
	sum := 0
	for _, it := range c.Items {
		sum += it.Price * it.Qty
	}
	return sum
}

// Add inserts a line, merging quantity when the same product+size is present.
func (c *Cart) Add(item Item) {
	for i := range c.Items {
		if c.Items[i].ProductUUID == item.ProductUUID && c.Items[i].ML == item.ML {
			c.Items[i].Qty += item.Qty
			return
		}
	}
	c.Items = append(c.Items, item)
}

// Remove drops the line matching the product+size.
func (c *Cart) Remove(productUUID string, ml int) {
	kept := c.Items[:0]
	for _, it := range c.Items {
		if it.ProductUUID == productUUID && it.ML == ml {
			continue
		}
		kept = append(kept, it)
	}
	c.Items = kept
}
